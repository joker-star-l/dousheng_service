package vo

import common "github.com/joker-star-l/dousheng_common/entity"

type LoginUser struct {
	Username string `json:"username" query:"username"`
	Password string `json:"password" query:"password"`
}

type UserLoginResponse struct {
	common.Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type UserInfoResponse struct {
	common.Response
	UserInfo
}
