package service

import (
	"context"
	common "github.com/joker-star-l/dousheng_common/entity"
	api "github.com/joker-star-l/dousheng_idls/user/kitex_gen/api"
)

type UserImpl struct{}

func (s *UserImpl) UserInfo(ctx context.Context, userId int64, queryId int64) (resp *api.UserInfoResponse, err error) {
	info, err := UserInfo(userId, queryId)
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
		Id:              info.Id,
		Name:            info.Name,
		FollowCount:     info.FollowCount,
		FollowerCount:   info.FollowerCount,
		IsFollow:        info.IsFollow,
		Avatar:          info.Avatar,
		BackgroundImage: info.BackgroundImage,
		Signature:       info.Signature,
		TotalFavorited:  info.TotalFavorited,
		WorkCount:       info.WorkCount,
		FavoriteCount:   info.FavoriteCount,
	}, nil
}
