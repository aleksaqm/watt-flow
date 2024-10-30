package handler

import "github.com/gin-gonic/gin"

type HealthHandler struct{}

func (h HealthHandler) Status(c *gin.Context) {
	c.String(200, "Alive")
}
