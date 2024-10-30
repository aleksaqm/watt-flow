package server

import (
	"watt-flow/handler"
	"watt-flow/service"
	"watt-flow/util"
)

type Server struct {
	Logger         util.Logger
	UserController handler.UserHandler
	userService    service.IUserService
}

func NewServer(logger util.Logger, userService service.IUserService, userController handler.UserHandler) *Server {
	return &Server{
		Logger:         logger,
		UserController: userController,
		userService:    userService,
	}
}
