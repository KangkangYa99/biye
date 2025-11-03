package redis

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Nil = redis.Nil

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "redis:6379" // 默认值
	}

	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")
	db := 0
	if dbStr != "" {
		if d, err := strconv.Atoi(dbStr); err == nil {
			db = d
		}
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis连接失败 :%v", err)
	} else {
		log.Println("✅ Redis连接成功:", addr)
	}
}
func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}
