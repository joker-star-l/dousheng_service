package service

import (
	"dousheng_service/user/domain/user/entity"
	"dousheng_service/user/infrastructure/config"
	"dousheng_service/user/infrastructure/db"
	"dousheng_service/user/interfaces/vo"
	"errors"
	"github.com/joker-star-l/dousheng_common/config/log"
	common "github.com/joker-star-l/dousheng_common/entity"
)

func Register(userVO *vo.LoginUser) (*common.TokenUser, error) {
	tx := db.DB.Where("name = ?", userVO.Username).Limit(1).Select("name").Find(&entity.User{})
	if tx.RowsAffected > 0 {
		return nil, errors.New("用户名已存在")
	}
	user := &entity.User{Name: userVO.Username, Password: userVO.Password}
	user.Id = config.GenerateId()
	tx = db.DB.Create(user)
	if tx.Error != nil {
		log.Slog.Errorf("create user error: %v", tx.Error.Error())
		return nil, errors.New("创建用户失败")
	}
	return &common.TokenUser{Id: user.Id, Name: user.Name}, nil
}

func Login(userVO *vo.LoginUser) (*common.TokenUser, error) {
	user := &entity.User{}
	tx := db.DB.Where("name = ? and password = ?", userVO.Username, userVO.Password).Limit(1).Find(user)
	if tx.RowsAffected < 1 {
		return nil, errors.New("用户名或密码错误")
	}
	return &common.TokenUser{Id: user.Id, Name: user.Name}, nil
}

func UserInfo(userId int64) (*vo.UserInfo, error) {
	user := &entity.User{}
	tx := db.DB.Limit(1).Find(user, userId)
	if tx.RowsAffected < 1 {
		return nil, errors.New("用户不存在")
	}
	return &vo.UserInfo{Id: user.Id, Name: user.Name}, nil
}
