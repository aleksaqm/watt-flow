package middleware

import (
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

type InvalidationRule struct {
	Pattern string
	Params  []string
}

func CacheInvalidationMiddleware(store persist.CacheStore, rules ...InvalidationRule) gin.HandlerFunc {
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
