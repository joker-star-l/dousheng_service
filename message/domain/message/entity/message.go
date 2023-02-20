package entity

import common "github.com/joker-star-l/dousheng_common/entity"

type Message struct {
	common.Model
	UserFrom int64  `json:"user_from" gorm:"column:user_from"`
	UserTo   int64  `json:"user_to" gorm:"column:user_to"`
	Content  string `json:"content" gorm:"column:content"`
}

func (m *Message) TableName() string {
	return "message"
}

//const RedisKeyMessage = "message"
