package entity

import (
	"dousheng_service/video/infrastructure/gorm"
	"dousheng_service/video/infrastructure/redis"
	"dousheng_service/video/infrastructure/snowflake"
	"errors"
	"fmt"
	"github.com/joker-star-l/dousheng_common/config/log"
	common "github.com/joker-star-l/dousheng_common/entity"
)

type VideoFavorite struct {
	common.Model
	UserId      int64 `json:"user_id" gorm:"column:user_id"`
	VideoId     int64 `json:"video_id" gorm:"column:video_id"`
	VideoUserId int64 `json:"video_user_id" gorm:"column:video_user_id"`
}

func (c *VideoFavorite) TableName() string {
	return "video_favorite"
}

var VideoFavoriteRepo = &VideoFavoriteRepository{}

type VideoFavoriteRepository struct{}

func (r *VideoFavoriteRepository) Create(videoFavorite *VideoFavorite) error {
	// 检查是否存在记录
	tx := gorm.DB.Select("id").Where("user_id = ? and video_id = ?", videoFavorite.UserId, videoFavorite.VideoId).Limit(1).Find(&VideoFavorite{})
	if tx.RowsAffected > 0 {
		return errors.New("不能重复点赞")
	}
	// save db
	videoFavorite.Id = snowflake.GenerateId()
	tx = gorm.DB.Create(videoFavorite)
	if tx.Error != nil {
		log.Slog.Errorf("create video_favorite error: %v", tx.Error.Error())
		return errors.New("点赞失败")
	}

	go func() {
		// 出错重试
		err := errors.New("start")
		for i := 0; i < 3 && err != nil; i++ {
			pipe := redis.Client.TxPipeline()
			// 更新 redis 视频点赞数
			pipe.HIncrBy(fmt.Sprintf("%s:%d", RedisKeyVideoStatistics, videoFavorite.VideoId), RedisHKeyVideoFavoriteCount, 1)
			// 更新 redis 用户喜欢数量
			pipe.HIncrBy(fmt.Sprintf("%s:%d", RedisKeyUserStatistics, videoFavorite.UserId), RedisHKeyUserFavoriteCount, 1)
			// 更新 redis 用户获赞数量
			pipe.HIncrBy(fmt.Sprintf("%s:%d", RedisKeyUserStatistics, videoFavorite.VideoUserId), RedisHKeyUserTotalFavorited, 1)
			_, err = pipe.Exec()
		}
		if err != nil {
			log.Slog.Errorln(err)
		}
	}()

	return nil
}

func (r *VideoFavoriteRepository) Delete(userId int64, videoId int64) error {
	videoFavorite := &VideoFavorite{}
	// 查询
	tx := gorm.DB.Where("user_id = ? and video_id = ?", userId, videoId).Select("id, video_user_id").Find(videoFavorite)
	if tx.RowsAffected > 0 {
		// 删除
		gorm.DB.Delete(&VideoFavorite{}, videoFavorite.Id)
		go func() {
			// 出错重试
			err := errors.New("start")
			for i := 0; i < 3 && err != nil; i++ {
				pipe := redis.Client.TxPipeline()
				// 更新 redis 视频点赞数
				pipe.HIncrBy(fmt.Sprintf("%s:%d", RedisKeyVideoStatistics, videoId), RedisHKeyVideoFavoriteCount, -1)
				// 更新 redis 用户喜欢数量
				pipe.HIncrBy(fmt.Sprintf("%s:%d", RedisKeyUserStatistics, userId), RedisHKeyUserFavoriteCount, -1)
				// 更新 redis 用户获赞数量
				pipe.HIncrBy(fmt.Sprintf("%s:%d", RedisKeyUserStatistics, videoFavorite.VideoUserId), RedisHKeyUserTotalFavorited, -1)
				_, err = pipe.Exec()
			}
			if err != nil {
				log.Slog.Errorln(err)
			}
		}()
	}
	return nil
}
