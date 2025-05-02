package utils

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var Red *redis.Client

func OpenRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

}
