package entity

import common "github.com/joker-star-l/dousheng_common/entity"

type VideoFavorite struct {
	common.Model
	UserId  int64 `json:"user_id" gorm:"column:user_id"`
	VideoId int64 `json:"video_id" gorm:"column:video_id"`
}

func (c *VideoFavorite) TableName() string {
	return "video_favorite"
}
