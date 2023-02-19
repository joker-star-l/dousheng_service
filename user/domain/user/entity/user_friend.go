package entity

import common "github.com/joker-star-l/dousheng_common/entity"

type UserFriend struct {
	common.Model
	User0 int64 `json:"user0" gorm:"column:user0"`
	User1 int64 `json:"user1" gorm:"column:user1"`
}

func (u *UserFriend) TableName() string {
	return "user_friend"
}
