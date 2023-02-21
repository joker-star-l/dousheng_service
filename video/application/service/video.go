package service

import (
	"context"
	"dousheng_service/video/config"
	"dousheng_service/video/domain/video/entity"
	"dousheng_service/video/infrastructure/gorm"
	"dousheng_service/video/infrastructure/kitex"
	my_minio "dousheng_service/video/infrastructure/minio"
	"dousheng_service/video/infrastructure/redis"
	"dousheng_service/video/infrastructure/snowflake"
	"dousheng_service/video/interfaces/vo"
	"errors"
	"fmt"
	"github.com/joker-star-l/dousheng_common/config/log"
	common "github.com/joker-star-l/dousheng_common/entity"
	util_redis "github.com/joker-star-l/dousheng_common/util/redis"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"strings"
	"time"
)

var VideoType = &[]string{".mp4", ".flv", ".f4v", ".webm"}

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

func Publish(userId int64, title string, file *multipart.FileHeader) error {
	reader, err := file.Open()
	if err != nil {
		log.Slog.Errorln(err)
		return errors.New("视频上传失败")
	}
	fileSuffix := ""
	for _, suffix := range *VideoType {
		if strings.HasSuffix(file.Filename, suffix) {
			fileSuffix = suffix
			break
		}
	}
	if fileSuffix == "" {
		return errors.New("视频格式错误")
	}
	videoId := snowflake.GenerateId()
	fileAddr := fmt.Sprintf("video/%d%s", videoId, fileSuffix)
	_, err = my_minio.Client.PutObject(context.Background(), config.C.Minio.Bucket, fileAddr, reader, file.Size, minio.PutObjectOptions{})
	if err != nil {
		log.Slog.Errorln(err)
		return errors.New("视频上传失败")
	}
	video := &entity.Video{
		Title:    title,
		PlayUrl:  my_minio.GetFullAddress(fileAddr),
		CoverUrl: "",
		UserId:   userId,
	}
	video.Id = videoId
	err = entity.VideoRepo.Create(video)
	if err != nil {
		return err
	}
	return nil
}

func Favorite(userId int64, videoId int64) error {
	// 检查视频是否存在
	video := &entity.Video{}
	tx := gorm.DB.Select("user_id").Limit(1).Find(video, videoId)
	if tx.RowsAffected < 1 {
		return errors.New("视频不存在")
	}
	videoFavorite := &entity.VideoFavorite{UserId: userId, VideoId: videoId, VideoUserId: video.UserId}
	return entity.VideoFavoriteRepo.Create(videoFavorite)
}

func CancelFavorite(userId int64, videoId int64) error {
	return entity.VideoFavoriteRepo.Delete(userId, videoId)
}

// GetFavoriteList userId: 用户自身id, queryId: 被查询的用户id
func GetFavoriteList(userId int64, queryId int64) ([]vo.VideoInfo, error) {
	// 查询列表
	var favoriteList []entity.VideoFavorite
	gorm.DB.Where("user_id", queryId).Select("video_id, video_user_id").Find(&favoriteList)
	result := make([]vo.VideoInfo, 0, len(favoriteList))
	for _, favorite := range favoriteList {
		// rpc
		userInfo := GetUserInfo(userId, favorite.VideoUserId)
		// 查询video
		video := &entity.Video{}
		gorm.DB.Limit(1).Find(video, favorite.VideoId)
		// 组装
		videoInfo := GetVideoInfo(userId, video, userInfo)
		result = append(result, *videoInfo)
	}
	return result, nil
}

func Comment(userId int64, videoId int64, commentText string) (*vo.Comment, error) {
	// 检查视频是否存在
	video := &entity.Video{}
	tx := gorm.DB.Select("user_id").Limit(1).Find(video, videoId)
	if tx.RowsAffected < 1 {
		return nil, errors.New("视频不存在")
	}
	// 存储
	comment := &entity.VideoComment{
		UserId:      userId,
		VideoId:     videoId,
		VideoUserId: video.UserId,
		Comment:     commentText,
	}
	err := entity.VideoCommentRepo.Create(comment)
	if err != nil {
		return nil, err
	}
	// RPC
	userInfo := GetUserInfo(userId, userId)
	// 组装
	return &vo.Comment{
		Id:         comment.Id,
		User:       *userInfo,
		Content:    comment.Comment,
		CreateDate: comment.CreatedAt.Format("01-02"),
	}, nil
}

func DeleteComment(userId int64, commentId int64) error {
	return entity.VideoCommentRepo.Delete(userId, commentId)
}

func GetCommentList(userId int64, videoId int64) ([]vo.Comment, error) {
	var commentList []entity.VideoComment
	gorm.DB.Where("video_id = ?", videoId).Order("created_at desc").Find(&commentList)
	result := make([]vo.Comment, 0, len(commentList))
	for _, comment := range commentList {
		// RPC
		userInfo := GetUserInfo(userId, comment.UserId)
		// 组装
		result = append(result, vo.Comment{
			Id:         comment.Id,
			User:       *userInfo,
			Content:    comment.Comment,
			CreateDate: comment.CreatedAt.Format("01-02"),
		})
	}
	return result, nil
}

// GetUserInfo userId: 用户自身id, queryId: 被查询的用户id
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
