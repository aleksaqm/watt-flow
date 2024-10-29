package server

import (
	"watt-flow/controller"
	"watt-flow/service"
	"watt-flow/util"
)

type Server struct {
	Logger         util.Logger
	UserController controller.UserController
	userService    service.IUserService
}

func NewServer(logger util.Logger, userService service.IUserService, userController controller.UserController) *Server {
	return &Server{
		Logger:         logger,
		UserController: userController,
		userService:    userService,
	}
}
