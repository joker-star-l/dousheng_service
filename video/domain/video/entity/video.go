package entity

import common "github.com/joker-star-l/dousheng_common/entity"

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
