package service

import (
	"context"
	common "github.com/joker-star-l/dousheng_common/entity"
	api "github.com/joker-star-l/dousheng_idls/user/kitex_gen/api"
)

type UserImpl struct{}

func (s *UserImpl) UserInfo(ctx context.Context, userId int64) (resp *api.UserInfoResponse, err error) {
	info, err := UserInfo(userId, userId)
	if err != nil {
		return &api.UserInfoResponse{
			Response: &api.Response{
				StatusCode: common.StatusError,
				StatusMsg:  err.Error(),
			},
		}, nil
	}
	return &api.UserInfoResponse{
		Response: &api.Response{
			StatusCode: common.StatusSuccess,
			StatusMsg:  "success",
		},
		Id:   info.Id,
		Name: info.Name,
	}, nil
}
