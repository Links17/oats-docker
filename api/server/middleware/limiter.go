package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"

	"oats-docker/api/server/httputils"
	"oats-docker/pkg/util/lru"
)

const (
	capacity = 100
	quantum  = 20
	cap      = 200
)

// UserRateLimiter 针对每个用户的请求进行限速
// TODO 限速大小从配置中读取
func UserRateLimiter() gin.HandlerFunc {
	// 初始化一个 LRU Cache
	cache, _ := lru.NewLRUCache(cap)

	return func(c *gin.Context) {
		r := httputils.NewResponse()
		// 把 key: clientIP value: *ratelimit.Bucket 存入 LRU Cache 中
		clientIP := c.ClientIP()
		if !cache.Contains(clientIP) {
			cache.Add(clientIP, ratelimit.NewBucketWithQuantum(time.Second, capacity, quantum))
			return
		}
		// 通过 ClientIP 取出 bucket
		val := cache.Get(clientIP)
		if val == nil {
			return
		}

		// 判断是否还有可用的 bucket
		bucket := val.(*ratelimit.Bucket)
		if bucket.TakeAvailable(1) == 0 {
			r.SetCode(http.StatusGatewayTimeout)
			httputils.SetFailed(c, r, fmt.Errorf("the system is busy. please try again later"))
			c.Abort()
			return
		}
	}
}

// Limiter TODO 总量限速
func Limiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("TODO")
	}
}
