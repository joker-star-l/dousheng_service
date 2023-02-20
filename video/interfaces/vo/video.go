package vo

import common "github.com/joker-star-l/dousheng_common/entity"

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

type VideoInfo struct {
	Id            int64    `json:"id"`
	Author        UserInfo `json:"author"`
	PlayUrl       string   `json:"play_url"`
	CoverUrl      string   `json:"cover_url"`
	FavoriteCount int64    `json:"favorite_count"`
	CommentCount  int64    `json:"comment_count"`
	IsFavorite    bool     `json:"is_favorite"`
	Title         string   `json:"title"`
	CreateTime    int64    `json:"create_time"`
}

type VideoInfoListResponse struct {
	common.Response
	VideoList []VideoInfo `json:"video_list"`
	NextTime  int64       `json:"next_time,omitempty"`
}

type Comment struct {
	Id         int64    `json:"id"`
	User       UserInfo `json:"user"`
	Content    string   `json:"content"`
	CreateDate string   `json:"create_date"`
}

type CommentResponse struct {
	common.Response
	Comment Comment `json:"comment"`
}

type CommentListResponse struct {
	common.Response
	CommentList []Comment `json:"comment_list"`
}
