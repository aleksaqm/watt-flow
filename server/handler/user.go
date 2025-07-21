package handler

import (
	"strconv"

	"watt-flow/dto"
	"watt-flow/service"
	"watt-flow/util"

	"github.com/gin-gonic/gin"
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

func (u *UserHandler) Suspend(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = u.service.Suspend(userId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": "success"})
}

func (u *UserHandler) SuspendClerk(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = u.service.SuspendClerk(userId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": "success"})
}

func (u *UserHandler) Unsuspend(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = u.service.Unsuspend(userId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": "success"})
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

func (h *UserHandler) Query(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "asc")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid page parameter"})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid pageSize parameter"})
		return
	}

	var searchParams dto.UserSearchParams
	if err := c.BindJSON(&searchParams); err != nil {
		c.JSON(400, gin.H{"error": "Invalid search parameter"})
		return
	}

	params := dto.UserQueryParams{
		Page:      pageInt,
		PageSize:  pageSizeInt,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Search:    searchParams,
	}
	h.logger.Info(params)
	users, total, err := h.service.Query(&params)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"users": users, "total": total})
}

func (u *UserHandler) RegisterClerk(c *gin.Context) {
	var clerk dto.ClerkRegisterDto
	if err := c.BindJSON(&clerk); err != nil {
		u.logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	data, err := u.service.RegisterClerk(&clerk)
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

func (u *UserHandler) ReturnOk(c *gin.Context) {
	c.JSON(200, gin.H{"data": "ok"})
}

func NewUserHandler(userService service.IUserService, logger util.Logger) *UserHandler {
	return &UserHandler{
		service: userService,
		logger:  logger,
	}
}
