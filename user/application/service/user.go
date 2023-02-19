package service

import (
	"context"
	"dousheng_service/user/domain/user/entity"
	"dousheng_service/user/infrastructure/config"
	"dousheng_service/user/infrastructure/gorm"
	my_minio "dousheng_service/user/infrastructure/minio"
	"dousheng_service/user/infrastructure/snowflake"
	"dousheng_service/user/interfaces/vo"
	"encoding/base64"
	"errors"
	"github.com/joker-star-l/dousheng_common/config/log"
	common "github.com/joker-star-l/dousheng_common/entity"
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

func UserInfo(userId int64) (*vo.UserInfo, error) {
	user := &entity.User{}
	tx := gorm.DB.Limit(1).Find(user, userId)
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

	// TODO count 信息
	return &vo.UserInfo{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     0,
		FollowerCount:   0,
		IsFollow:        false,
		Avatar:          avatar,
		BackgroundImage: backgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  0,
		WorkCount:       0,
		FavoriteCount:   0,
	}, nil
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
