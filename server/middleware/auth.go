package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"watt-flow/service"
	"watt-flow/util"
)

type AuthMiddleware struct {
	service *service.AuthService
	logger  util.Logger
}

func NewAuthMiddleware(service *service.AuthService, logger util.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		service: service,
		logger:  logger,
	}
}

func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, err := m.service.Authorize(authToken)
			if authorized {
				c.Next()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			m.logger.Error(err)
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "you are not authorized",
		})
		c.Abort()
	}
}
