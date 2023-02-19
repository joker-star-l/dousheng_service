package vo

import common "github.com/joker-star-l/dousheng_common/entity"

type LoginUser struct {
	Username string `json:"username" query:"username" vd:"len($)>0; msg:'用户名长度必须大于0'"`
	Password string `json:"password" query:"password" vd:"len($)>7; msg:'密码长度必须大于等于8'"`
}

type UserLoginResponse struct {
	common.Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfo struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	FollowCount     int64  `json:"follow_count"`
	FollowerCount   int64  `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  int64  `json:"total_favorited"`
	WorkCount       int64  `json:"work_count"`
	FavoriteCount   int64  `json:"favorite_count"`
}

type UserInfoResponse struct {
	common.Response
	User UserInfo `json:"user"`
}

type UserInfoListResponse struct {
	common.Response
	UserList []UserInfo `json:"user_list"`
}
