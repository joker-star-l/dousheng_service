package entity

import common "github.com/joker-star-l/dousheng_common/entity"

type User struct {
	common.Model
	Name     string `json:"name" gorm:"column:name"`
	Password string `json:"password" gorm:"column:password"`
}

func (u *User) TableName() string {
	return "user"
}
