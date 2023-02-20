package service

import (
	"context"
	"dousheng_service/video/domain/video/entity"
	"dousheng_service/video/infrastructure/gorm"
	"dousheng_service/video/infrastructure/kitex"
	"dousheng_service/video/infrastructure/redis"
	"dousheng_service/video/interfaces/vo"
	"fmt"
	"github.com/joker-star-l/dousheng_common/config/log"
	common "github.com/joker-star-l/dousheng_common/entity"
	util_redis "github.com/joker-star-l/dousheng_common/util/redis"
	"time"
)

func Feed(userId int64, latestTime time.Time) ([]vo.VideoInfo, error) {
	var videoList []entity.Video
	gorm.DB.Where("created_at < ?", latestTime).Order("created_at desc").Limit(30).Find(&videoList)
	var result = make([]vo.VideoInfo, 0, len(videoList))
	for _, video := range videoList {
		// RPC调用
		userInfo := GetUserInfo(userId, video.UserId)
		// 组装
		videoInfo := GetVideoInfo(userId, &video, userInfo)
		// 添加到返回列表
		result = append(result, *videoInfo)
	}
	return result, nil
}

func GetPublishList(userId int64, queryId int64) ([]vo.VideoInfo, error) {
	var videoList []entity.Video
	gorm.DB.Where("user_id = ?", queryId).Order("created_at desc").Find(&videoList)
	var result = make([]vo.VideoInfo, 0, len(videoList))
	// RPC调用
	userInfo := GetUserInfo(userId, queryId)
	for _, video := range videoList {
		// 组装
		videoInfo := GetVideoInfo(userId, &video, userInfo)
		// 添加到返回列表
		result = append(result, *videoInfo)
	}
	return result, nil
}

func GetUserInfo(userId int64, queryId int64) *vo.UserInfo {
	userInfo := &vo.UserInfo{}
	info, _ := kitex.UserClient.UserInfo(context.Background(), userId, queryId)
	if info == nil {
		log.Slog.Errorln("RPC 调用失败")
	} else if info.Response.StatusCode != common.StatusSuccess {
		log.Slog.Errorln(info.Response.StatusMsg)
	} else {
		userInfo = &vo.UserInfo{
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
		}
	}
	return userInfo
}

func GetVideoInfo(userId int64, video *entity.Video, userInfo *vo.UserInfo) *vo.VideoInfo {
	// 组装
	videoInfo := &vo.VideoInfo{
		Id:         video.Id,
		Author:     *userInfo,
		PlayUrl:    video.PlayUrl,
		CoverUrl:   video.CoverUrl,
		Title:      video.Title,
		CreateTime: video.CreatedAt.Unix() * 1000,
	}
	// 获取 count
	count, err := redis.Client.HGetAll(fmt.Sprintf("%s:%d", entity.RedisKeyVideoStatistics, video.Id)).Result()
	if err != nil {
		log.Slog.Errorln(err)
	} else {
		videoInfo.FavoriteCount = util_redis.ParseCount(count[entity.RedisHKeyVideoFavoriteCount])
		videoInfo.CommentCount = util_redis.ParseCount(count[entity.RedisHKeyVideoCommentCount])
	}
	// 是否喜欢
	if userId != 0 {
		tx := gorm.DB.Where("user_id = ? and video_id = ?", userId, video.Id).Select("id").Limit(1).Find(&entity.VideoFavorite{})
		if tx.RowsAffected > 0 {
			videoInfo.IsFavorite = true
		}
	}
	return videoInfo
}
