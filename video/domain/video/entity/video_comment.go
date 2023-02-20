package entity

import common "github.com/joker-star-l/dousheng_common/entity"

type VideoComment struct {
	common.Model
	UserId  int64  `json:"user_id" gorm:"column:user_id"`
	VideoId int64  `json:"video_id" gorm:"column:video_id"`
	Comment string `json:"comment" gorm:"column:comment"`
}

func (c *VideoComment) TableName() string {
	return "video_comment"
}
