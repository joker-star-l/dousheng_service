package entity

import (
	"dousheng_service/user/infrastructure/gorm"
	"dousheng_service/user/infrastructure/redis"
	"dousheng_service/user/infrastructure/snowflake"
	"errors"
	"fmt"
	"github.com/joker-star-l/dousheng_common/config/log"
	common "github.com/joker-star-l/dousheng_common/entity"
	gorm2 "gorm.io/gorm"
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
	// 检查是否是双向关注
	both := false
	tx = gorm.DB.Select("id").Where("user_from = ? and user_to = ?", userFollow.UserTo, userFollow.UserFrom).Limit(1).Find(&UserFollow{})
	if tx.RowsAffected > 0 {
		both = true
	}
	// save db
	err := gorm.DB.Transaction(func(tx *gorm2.DB) error {
		// 关注
		userFollow.Id = snowflake.GenerateId()
		if tx.Create(userFollow).Error != nil {
			log.Slog.Errorf("create user_follow error: %v", tx.Error.Error())
			return errors.New("关注用户失败")
		}
		// 成为好友
		if both {
			userFriend := &UserFriend{User0: userFollow.UserFrom, User1: userFollow.UserTo}
			userFriend.Id = snowflake.GenerateId()
			if tx.Create(userFriend).Error != nil {
				log.Slog.Errorf("create user_friend error: %v", tx.Error.Error())
				return errors.New("成为好友失败")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	go func() {
		pipe := redis.Client.TxPipeline()
		// 更新 redis 关注数
		pipe.HIncrBy(fmt.Sprintf("%s:%d", RedisKeyUserStatistics, userFollow.UserFrom), RedisHKeyUserFollowCount, 1)
		// 更新 redis 粉丝数
		pipe.HIncrBy(fmt.Sprintf("%s:%d", RedisKeyUserStatistics, userFollow.UserTo), RedisHKeyUserFollowerCount, 1)
		_, err := pipe.Exec()
		// 出错重试一次
		if err != nil {
			log.Slog.Errorln(err)
			pipe.Exec()
		}
	}()

	return nil
}

func (r *UserFollowRepository) Delete(from int64, to int64) error {
	gorm.DB.Transaction(func(tx *gorm2.DB) error {
		if tx.Where("user_from = ? and user_to = ?", from, to).Delete(&UserFollow{}).RowsAffected > 0 {
			tx.Where("(user0 = ? and user1 = ?) or (user0 = ? and user1 = ?)", from, to, to, from).Delete(&UserFriend{})
			go func() {
				pipe := redis.Client.TxPipeline()
				// 更新 redis 关注数
				pipe.HIncrBy(fmt.Sprintf("%s:%d", RedisKeyUserStatistics, from), RedisHKeyUserFollowCount, -1)
				// 更新 redis 粉丝数
				pipe.HIncrBy(fmt.Sprintf("%s:%d", RedisKeyUserStatistics, to), RedisHKeyUserFollowerCount, -1)
				_, err := pipe.Exec()
				// 出错重试一次
				if err != nil {
					log.Slog.Errorln(err)
					pipe.Exec()
				}
			}()
		}
		return nil
	})
	return nil
}
