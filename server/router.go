package server

import (
	"watt-flow/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controller.HealthController)
	router.GET("/health", health.Status)

	v1 := router.Group("v1")
	{
		authGroup := v1.Group("auth")
		{
			authGroup.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})
		}
	}

	return router
}
