package service

import (
	"context"
	"dousheng_service/message/domain/message/entity"
	"dousheng_service/message/infrastructure/gorm"
	common "github.com/joker-star-l/dousheng_common/entity"
	"github.com/joker-star-l/dousheng_idls/message/kitex_gen/api"
)

type MessageImpl struct{}

func (s *MessageImpl) LatestMessage(ctx context.Context, userId int64, friendId int64) (*api.LatestMessageResponse, error) {
	message := &entity.Message{}
	gorm.DB.Where("(user_from = ? and user_to = ?) or (user_from = ? and user_to = ?)", userId, friendId, friendId, userId).
		Order("created_at desc").Limit(1).Find(message)

	return &api.LatestMessageResponse{
		Response: &api.Response{StatusCode: common.StatusSuccess, StatusMsg: "success"},
		Message:  message.Content,
		From:     message.UserFrom,
		To:       message.UserTo,
	}, nil
}
