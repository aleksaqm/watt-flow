package middleware

import (
	"context"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"strings"
)

type InvalidationRule struct {
	Pattern string
	Params  []string
}

func CacheInvalidationMiddleware(store persist.CacheStore, redisClient *redis.Client, rules ...InvalidationRule) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			for _, rule := range rules {
				cacheKey := rule.Pattern
				paramsMatch := true
				for _, paramName := range rule.Params {
					paramValue := c.Param(paramName)
					if paramValue == "" {
						paramsMatch = false
						break
					}
					cacheKey = strings.Replace(cacheKey, ":"+paramName, paramValue, 1)
				}

				if paramsMatch {
					if strings.Contains(cacheKey, "*") {
						go func(pattern string) {
							ctx := context.Background()
							iter := redisClient.Scan(ctx, 0, pattern, 0).Iterator()
							keysToDelete := []string{}
							for iter.Next(ctx) {
								keysToDelete = append(keysToDelete, iter.Val())
							}
							if err := iter.Err(); err != nil {
								log.Printf("ERROR: Redis SCAN failed for pattern '%s': %v", pattern, err)
								return
							}
							if len(keysToDelete) > 0 {
								redisClient.Del(ctx, keysToDelete...)
								log.Printf("INFO: Cache invalidated for %d keys matching pattern: %s", len(keysToDelete), pattern)
							}
						}(cacheKey)

					} else {
						err := store.Delete(cacheKey)
						if err != nil {
							log.Printf("ERROR: Failed to invalidate cache for key '%s': %v", cacheKey, err)
						} else {
							log.Printf("INFO: Cache invalidated for key: %s", cacheKey)
						}
					}
				}
			}
		}
	}
}
