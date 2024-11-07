package handler

import (
	"github.com/gin-gonic/gin"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/service"
	"watt-flow/util"
)

type UserHandler struct {
	service service.IUserService
	logger  util.Logger
}

func (u UserHandler) Login(c *gin.Context) {
	var loginCredentials dto.LoginDto
	if err := c.BindJSON(&loginCredentials); err != nil {
		u.logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
	}
	token, err := u.service.Login(loginCredentials)
	if err != nil {
		u.logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"token": token})
	}
}

func (u UserHandler) Register(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		u.logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
	}
	data, _ := u.service.Create(&user)
	c.JSON(200, gin.H{"data": data})
}

func (u UserHandler) GetById(c *gin.Context) {
	data, _ := u.service.FindById(1)
	u.logger.Info("radi handler")
	c.JSON(200, gin.H{"data": data})
}

func (u UserHandler) Create(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		u.logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
	}
	data, _ := u.service.Create(&user)
	c.JSON(200, gin.H{"data": data})
}

func NewUserHandler(userService service.IUserService, logger util.Logger) *UserHandler {
	return &UserHandler{
		service: userService,
		logger:  logger,
	}
}
