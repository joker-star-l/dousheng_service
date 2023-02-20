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

type VideoComment struct {
	common.Model
	UserId      int64  `json:"user_id" gorm:"column:user_id"`
	VideoId     int64  `json:"video_id" gorm:"column:video_id"`
	VideoUserId int64  `json:"video_user_id" gorm:"column:video_user_id"`
	Comment     string `json:"comment" gorm:"column:comment"`
}

func (c *VideoComment) TableName() string {
	return "video_comment"
}

var VideoCommentRepo = &VideoCommentRepository{}

type VideoCommentRepository struct{}

func (r *VideoCommentRepository) Create(comment *VideoComment) error {
	comment.Id = snowflake.GenerateId()
	tx := gorm.DB.Create(comment)
	if tx.Error != nil {
		log.Slog.Errorf("create video_comment error: %v", tx.Error.Error())
		return errors.New("评论失败")
	}

	go func() {
		// 出错重试
		err := errors.New("start")
		for i := 0; i < 3 && err != nil; i++ {
			cmd := redis.Client.HIncrBy(fmt.Sprintf("%s:%d", RedisKeyVideoStatistics, comment.VideoId), RedisHKeyVideoCommentCount, 1)
			err = cmd.Err()
		}
		if err != nil {
			log.Slog.Errorln(err)
		}
	}()

	return nil
}

func (r VideoCommentRepository) Delete(userId int64, commentId int64) error {
	// 查询
	comment := &VideoComment{}
	tx := gorm.DB.Where("user_id = ? and id = ?", userId, commentId).Select("id, video_id").Limit(1).Find(comment)
	if tx.RowsAffected > 0 {
		gorm.DB.Delete(&VideoComment{}, comment.Id)
		go func() {
			err := errors.New("start")
			for i := 0; i < 3 && err != nil; i++ {
				cmd := redis.Client.HIncrBy(fmt.Sprintf("%s:%d", RedisKeyVideoStatistics, comment.VideoId), RedisHKeyVideoCommentCount, -1)
				err = cmd.Err()
			}
			if err != nil {
				log.Slog.Errorln(err)
			}
		}()
	}
	return nil
}
