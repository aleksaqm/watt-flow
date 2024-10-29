package controller

import (
	"github.com/gin-gonic/gin"
	"watt-flow/service"
	"watt-flow/util"
)

type UserController struct {
	service service.IUserService
	logger  util.Logger
}

func (u UserController) GetById(c *gin.Context) {
	data, _ := u.service.FindById("kdjfk")
	u.logger.Info("radi controller")
	c.JSON(200, gin.H{"data": data})
}

func NewUserController(userService service.IUserService, logger util.Logger) UserController {
	return UserController{
		service: userService,
		logger:  logger,
	}
}
