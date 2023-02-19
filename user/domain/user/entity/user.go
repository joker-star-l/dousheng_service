package entity

import common "github.com/joker-star-l/dousheng_common/entity"

type User struct {
	common.Model
	Name            string `json:"name" gorm:"column:name"`
	Password        string `json:"password" gorm:"column:password"`
	Avatar          string `json:"avatar" gorm:"column:avatar"`
	BackgroundImage string `json:"background_image" gorm:"column:background_image"`
	Signature       string `json:"signature" gorm:"column:signature"`
}

func (u *User) TableName() string {
	return "user"
}
