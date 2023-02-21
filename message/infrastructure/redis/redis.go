package redis

import (
	"dousheng_service/message/config"
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
