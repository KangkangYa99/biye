package middleware

import (
	"biye/share/redis"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimitConfig struct {
	KeyPrefix string        //key前缀
	Limit     int64         //限制次数
	Time      time.Duration //时间
}

func RateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := fmt.Sprintf("%s:%s", config.KeyPrefix, clientIP)
		ctx := context.Background()

		count, err := redis.RedisClient.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}
		if count == 1 {
			redis.RedisClient.Expire(ctx, key, config.Time)
		}
		if count > config.Limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"code":    429,
				"message": fmt.Sprintf("请求过于频繁,每分钟最多允许 %d 次请求，当前已请求 %d 次", config.Limit, count),
			})
			c.Abort()
			return
		}
		fmt.Printf("IP %s 在 %s 下的请求次数: %d/%d\n", clientIP, config.KeyPrefix, count, config.Limit)
		c.Next()
	}
}
