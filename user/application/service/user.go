package service

import (
	"context"
	"dousheng_service/user/domain/user/entity"
	"dousheng_service/user/infrastructure/config"
	"dousheng_service/user/infrastructure/gorm"
	my_minio "dousheng_service/user/infrastructure/minio"
	"dousheng_service/user/infrastructure/redis"
	"dousheng_service/user/infrastructure/snowflake"
	"dousheng_service/user/interfaces/vo"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/joker-star-l/dousheng_common/config/log"
	common "github.com/joker-star-l/dousheng_common/entity"
	util_redis "github.com/joker-star-l/dousheng_common/util/redis"
	"github.com/minio/minio-go/v7"
	"golang.org/x/crypto/bcrypt"
	"io"
	"strconv"
)

func Register(userVO *vo.LoginUser) (*common.TokenUser, error) {
	tx := gorm.DB.Where("name = ?", userVO.Username).Limit(1).Select("name").Find(&entity.User{})
	if tx.RowsAffected > 0 {
		return nil, errors.New("用户名已存在")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(userVO.Password), bcrypt.DefaultCost)
	user := &entity.User{
		Name:            userVO.Username,
		Password:        string(hash),
		Avatar:          my_minio.DefaultAvatarAddress,
		BackgroundImage: my_minio.DefaultBackgroundAddress,
		Signature:       "无",
	}
	user.Id = snowflake.GenerateId()
	tx = gorm.DB.Create(user)
	if tx.Error != nil {
		log.Slog.Errorf("create user error: %v", tx.Error.Error())
		return nil, errors.New("创建用户失败")
	}
	return &common.TokenUser{Id: strconv.FormatInt(user.Id, 10), Name: user.Name}, nil
}

func Login(userVO *vo.LoginUser) (*common.TokenUser, error) {
	user := &entity.User{}
	tx := gorm.DB.Where("name = ?", userVO.Username).Select("id, name, password").Limit(1).Find(user)
	if tx.RowsAffected < 1 {
		return nil, errors.New("用户名或密码错误")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userVO.Password))
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	return &common.TokenUser{Id: strconv.FormatInt(user.Id, 10), Name: user.Name}, nil
}

// UserInfo userId: 用户自身id, queryId: 被查询的用户id
func UserInfo(userId int64, queryId int64) (*vo.UserInfo, error) {
	user := &entity.User{}
	tx := gorm.DB.Limit(1).Find(user, queryId)
	if tx.RowsAffected < 1 {
		return nil, errors.New("用户不存在")
	}

	var avatar = ""
	reader, err := my_minio.Client.GetObject(context.Background(), config.C.Minio.Bucket, user.Avatar, minio.GetObjectOptions{})
	if err != nil {
		log.Slog.Errorln(err)
	} else {
		defer reader.Close()
		stat, err := reader.Stat()
		if err != nil {
			log.Slog.Errorln(err)
		} else {
			src := make([]byte, stat.Size)
			n, err := reader.Read(src)
			if n <= 0 || err != io.EOF {
				log.Slog.Errorln(err)
			} else {
				avatar = "data:image;base64," + base64.StdEncoding.EncodeToString(src)
			}
		}
	}

	var backgroundImage = ""
	reader, err = my_minio.Client.GetObject(context.Background(), config.C.Minio.Bucket, user.BackgroundImage, minio.GetObjectOptions{})
	if err != nil {
		log.Slog.Errorln(err)
	} else {
		defer reader.Close()
		stat, err := reader.Stat()
		if err != nil {
			log.Slog.Errorln(err)
		} else {
			src := make([]byte, stat.Size)
			n, err := reader.Read(src)
			if n <= 0 || err != io.EOF {
				log.Slog.Errorln(err)
			} else {
				backgroundImage = "data:image;base64," + base64.StdEncoding.EncodeToString(src)
			}
		}
	}

	userInfo := &vo.UserInfo{
		Id:              user.Id,
		Name:            user.Name,
		Avatar:          avatar,
		BackgroundImage: backgroundImage,
		Signature:       user.Signature,
	}

	// follow 信息
	if userId != 0 && userId != queryId {
		tx = gorm.DB.Select("id").Where("user_from = ? and user_to = ?", userId, queryId).Limit(1).Find(&entity.UserFollow{})
		if tx.RowsAffected > 0 {
			userInfo.IsFollow = true
		}
	}

	// count 信息
	result, err := redis.Client.HGetAll(fmt.Sprintf("%s:%d", entity.RedisKeyUserStatistics, queryId)).Result()
	if err != nil {
		log.Slog.Errorln(err)
	} else {
		userInfo.FollowCount = util_redis.ParseCount(result[entity.RedisHKeyUserFollowCount])
		userInfo.FollowerCount = util_redis.ParseCount(result[entity.RedisHKeyUserFollowerCount])
		userInfo.TotalFavorited = util_redis.ParseCount(result[entity.RedisHKeyUserTotalFavorited])
		userInfo.WorkCount = util_redis.ParseCount(result[entity.RedisHKeyUserWorkCount])
		userInfo.FavoriteCount = util_redis.ParseCount(result[entity.RedisHKeyUserFavoriteCount])
	}

	return userInfo, nil
}

func Follow(from int64, to int64) error {
	// 检查用户是否存在
	tx := gorm.DB.Select("id").Limit(1).Find(&entity.User{}, to)
	if tx.RowsAffected < 1 {
		return errors.New("用户不存在")
	}
	// 关注
	userFollow := &entity.UserFollow{UserFrom: from, UserTo: to}
	err := entity.UserFollowRepo.Create(userFollow)
	if err != nil {
		return err
	}
	return nil
}

func CancelFollow(from int64, to int64) error {
	return entity.UserFollowRepo.Delete(from, to)
}

// GetFollowList userId: 用户自身id, queryId: 被查询的用户id
func GetFollowList(userId int64, queryId int64) ([]vo.UserInfo, error) {
	var userFollowList []entity.UserFollow
	gorm.DB.Where("user_from = ?", queryId).Select("user_to").Find(&userFollowList)
	result := make([]vo.UserInfo, 0, len(userFollowList))
	if userId == queryId {
		for _, follow := range userFollowList {
			info, err := UserInfo(0, follow.UserTo)
			if err == nil {
				info.IsFollow = true
				result = append(result, *info)
			}
		}
	} else {
		for _, follow := range userFollowList {
			info, err := UserInfo(userId, follow.UserTo)
			if err == nil {
				result = append(result, *info)
			}
		}
	}
	return result, nil
}

// GetFollowerList userId: 用户自身id, queryId: 被查询的用户id
func GetFollowerList(userId int64, queryId int64) ([]vo.UserInfo, error) {
	var userFollowerList []entity.UserFollow
	gorm.DB.Where("user_to = ?", queryId).Select("user_from").Find(&userFollowerList)
	result := make([]vo.UserInfo, 0, len(userFollowerList))
	for _, follow := range userFollowerList {
		info, err := UserInfo(userId, follow.UserFrom)
		if err == nil {
			result = append(result, *info)
		}
	}
	return result, nil
}

func GetFriendList(userId int64) ([]vo.FriendInfo, error) {
	var userFriendList []entity.UserFriend
	gorm.DB.Where("user0 = ? or user1 = ?", userId, userId).Select("user0, user1").Find(&userFriendList)
	result := make([]vo.FriendInfo, 0, len(userFriendList))
	for _, friend := range userFriendList {
		var friendId int64
		if friend.User0 == userId {
			friendId = friend.User1
		} else {
			friendId = friend.User0
		}
		info, err := UserInfo(userId, friendId)
		if err == nil {
			// TODO RPC 调用
			result = append(result, vo.FriendInfo{
				UserInfo: *info,
				Message:  "FAKE",
				MsgType:  1,
			})
		}
	}
	return result, nil
}
