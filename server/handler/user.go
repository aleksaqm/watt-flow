package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"watt-flow/dto"
	"watt-flow/service"
	"watt-flow/util"
)

type UserHandler struct {
	service service.IUserService
	logger  util.Logger
}

func (u *UserHandler) Login(c *gin.Context) {
	var loginCredentials dto.LoginDto
	if err := c.BindJSON(&loginCredentials); err != nil {
		u.logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := u.service.Login(loginCredentials)
	if err != nil {
		u.logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"token": token})
	}
}

func (u *UserHandler) Register(c *gin.Context) {
	var user dto.RegistrationDto
	if err := c.BindJSON(&user); err != nil {
		u.logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	data, err := u.service.Register(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func (u *UserHandler) ActivateAccount(c *gin.Context) {
	token := c.Param("token")
	err := u.service.ActivateAccount(token)
	loginLink := "http://localhost:5173/"
	if err != nil {
		c.Data(200, "text/html; charset=utf-8", []byte(util.GenerateFailedActivationEmailBody(loginLink)))
		return
	}
	c.Data(200, "text/html; charset=utf-8", []byte(util.GenerateSuccessfulActivationEmailBody(loginLink)))
}

func (u *UserHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	data, _ := u.service.FindById(userId)
	u.logger.Info("radi handler")
	c.JSON(200, gin.H{"data": data})
}

func (u *UserHandler) Create(c *gin.Context) {
	var user dto.UserCreateDto
	if err := c.BindJSON(&user); err != nil {
		u.logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	data, err := u.service.Create(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"data": data})
	}
}

func (u *UserHandler) ChangeAdminPassword(c *gin.Context) {
	var passwords dto.NewPasswordDto
	if err := c.BindJSON(&passwords); err != nil {
		u.logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := u.service.ChangeAdminPassword(passwords)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"data": "Password changed successfully"})
	}
}

func (u *UserHandler) IsAdminActive(c *gin.Context) {
	isActive := u.service.IsAdminActive()
	c.JSON(200, gin.H{"active": isActive})
}

func (u *UserHandler) FindAdmins(c *gin.Context) {
	data, err := u.service.FindAllByRole("Admin")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	c.JSON(200, data)
}

func NewUserHandler(userService service.IUserService, logger util.Logger) *UserHandler {
	return &UserHandler{
		service: userService,
		logger:  logger,
	}
}
