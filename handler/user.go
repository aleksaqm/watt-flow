package handler

import (
	"github.com/gin-gonic/gin"
	"watt-flow/service"
	"watt-flow/util"
)

type UserHandler struct {
	service service.IUserService
	logger  util.Logger
}

func (u UserHandler) GetById(c *gin.Context) {
	data, _ := u.service.FindById("kdjfk")
	u.logger.Info("radi handler")
	c.JSON(200, gin.H{"data": data})
}

func NewUserHandler(userService service.IUserService, logger util.Logger) *UserHandler {
	return &UserHandler{
		service: userService,
		logger:  logger,
	}
}
