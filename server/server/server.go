package server

import (
	"watt-flow/db"
	"watt-flow/handler"
	"watt-flow/service"
	"watt-flow/util"
)

type Server struct {
	Logger      util.Logger
	UserHandler *handler.UserHandler
	userService service.IUserService
	AuthService *service.AuthService
	Db          db.Database
}

func NewServer(logger util.Logger, userService service.IUserService, authService *service.AuthService, userHandler *handler.UserHandler, db db.Database) *Server {
	return &Server{
		Logger:      logger,
		UserHandler: userHandler,
		userService: userService,
		AuthService: authService,
		Db:          db,
	}
}
