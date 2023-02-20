package vo

import common "github.com/joker-star-l/dousheng_common/entity"

type Message struct {
	Id         int64  `json:"id"`
	ToUserId   int64  `json:"to_user_id"`
	FromUserId int64  `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

type MessageListResponse struct {
	common.Response
	MessageList []Message `json:"message_list"`
}
