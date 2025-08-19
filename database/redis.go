package database

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // 你的 Redis 地址
		DB:   0,
	})
	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		panic("Redis 连接失败: " + err.Error())
	}
}
