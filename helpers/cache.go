package helpers

import (
	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     Config.RedisAddr,
		Password: Config.RedisPassword,
		DB:       0,
	})
}
