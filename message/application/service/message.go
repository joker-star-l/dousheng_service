package service

import (
	"dousheng_service/message/domain/message/entity"
	"dousheng_service/message/infrastructure/gorm"
	"dousheng_service/message/infrastructure/snowflake"
	"dousheng_service/message/interfaces/vo"
	"errors"
	"time"
)

func SendMessage(from int64, to int64, content string) error {
	//messageKey := fmt.Sprintf("%s:%d:%d", entity.RedisKeyMessage, from, to)
	//push := redis.Client.LPush(messageKey, content)
	//if push.Err() != nil {
	//	return errors.New("发送消息失败")
	//}
	message := &entity.Message{UserFrom: from, UserTo: to, Content: content}
	message.Id = snowflake.GenerateId()
	tx := gorm.DB.Create(message)
	if tx.Error != nil {
		return errors.New("消息发送失败")
	}
	return nil
}

func GetMessageList(from int64, to int64, lastTime int64) ([]vo.Message, error) {
	var messageList []entity.Message
	unix := time.Unix(lastTime/1000, 0)
	gorm.DB.Where("(user_from = ? and user_to = ?) or (user_from = ? and user_to = ?)", from, to, to, from).
		Where("created_at > ?", unix).Find(&messageList)
	result := make([]vo.Message, 0, len(messageList))
	for _, message := range messageList {
		result = append(result, vo.Message{
			Id:         message.Id,
			ToUserId:   message.UserTo,
			FromUserId: message.UserFrom,
			Content:    message.Content,
			CreateTime: message.CreatedAt.Unix() * 1000,
		})
	}
	return result, nil
}
