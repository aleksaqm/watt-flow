package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
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
			authorized, claims, err := m.service.Authorize(authToken)
			if authorized {
				c.Set("claims", claims)
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
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

func (m *AuthMiddleware) RoleMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("claims").(jwt.MapClaims)
		role := claims["role"].(string)
		if slices.Contains(roles, role) {
			c.Next()
			return
		}
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you don't have permission",
		})
		c.Abort()
	}
}
