package redis

import (
	"dousheng_service/user/infrastructure/config"
	"github.com/go-redis/redis"
)

var Client *redis.Client

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.C.Redis.Address,
		Password: config.C.Redis.Password,
		DB:       0,
	})
}
