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

type Video struct {
	common.Model
	Title    string `json:"title" gorm:"column:title"`
	PlayUrl  string `json:"play_url" gorm:"column:play_url"`
	CoverUrl string `json:"cover_url" gorm:"column:cover_url"`
	UserId   int64  `json:"user_id" gorm:"column:user_id"`
}

func (v *Video) TableName() string {
	return "video"
}

const (
	RedisKeyUserStatistics      = "user_statistics"
	RedisHKeyUserTotalFavorited = "total_favorited"
	RedisHKeyUserWorkCount      = "work_count"
	RedisHKeyUserFavoriteCount  = "favorite_count"
)

const (
	RedisKeyVideoStatistics     = "video_statistics"
	RedisHKeyVideoFavoriteCount = "favorite_count"
	RedisHKeyVideoCommentCount  = "comment_count"
)

var VideoRepo = &VideoRepository{}

type VideoRepository struct{}

func (r *VideoRepository) Create(video *Video) error {
	if video.Id == 0 {
		video.Id = snowflake.GenerateId()
	}
	tx := gorm.DB.Create(video)
	if tx.Error != nil {
		log.Slog.Errorf("create video error: %v", tx.Error.Error())
		return errors.New("创建视频失败")
	}

	go func() {
		// 出错重试
		err := errors.New("start")
		for i := 0; i < 3 && err != nil; i++ {
			cmd := redis.Client.HIncrBy(fmt.Sprintf("%s:%d", RedisKeyUserStatistics, video.UserId), RedisHKeyUserWorkCount, 1)
			err = cmd.Err()
		}
		if err != nil {
			log.Slog.Errorln(err)
		}
	}()

	return nil
}
