package entity

import (
	"dousheng_service/user/infrastructure/gorm"
	"dousheng_service/user/infrastructure/snowflake"
	"errors"
	"github.com/joker-star-l/dousheng_common/config/log"
	common "github.com/joker-star-l/dousheng_common/entity"
)

type UserFollow struct {
	common.Model
	UserFrom int64 `json:"user_from" gorm:"column:user_from"`
	UserTo   int64 `json:"user_to" gorm:"column:user_to"`
}

func (u *UserFollow) TableName() string {
	return "user_follow"
}

var UserFollowRepo = &UserFollowRepository{}

type UserFollowRepository struct{}

func (r *UserFollowRepository) Create(userFollow *UserFollow) error {
	// 检查是否存在记录
	tx := gorm.DB.Select("id").Where("user_from = ? and user_to = ?", userFollow.UserFrom, userFollow.UserTo).Limit(1).Find(&UserFollow{})
	if tx.RowsAffected > 0 {
		return errors.New("不能重复关注")
	}
	// save db
	userFollow.Id = snowflake.GenerateId()
	tx = gorm.DB.Create(userFollow)
	if tx.Error != nil {
		log.Slog.Errorf("create user_follow error: %v", tx.Error.Error())
		return errors.New("关注用户失败")
	}
	// TODO 更新 redis 关注数
	// TODO 更新 redis 粉丝数
	return nil
}

func (r *UserFollowRepository) Delete(from int64, to int64) error {
	tx := gorm.DB.Where("user_from = ? and user_to = ?", from, to).Delete(&UserFollow{})
	if tx.RowsAffected > 0 {
		// TODO 更新 redis 关注数
		// TODO 更新 redis 粉丝数
	}
	return nil
}
