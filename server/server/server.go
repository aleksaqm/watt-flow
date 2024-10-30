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
	Db          db.Database
}

func NewServer(logger util.Logger, userService service.IUserService, userHandler *handler.UserHandler, db db.Database) *Server {
	return &Server{
		Logger:      logger,
		UserHandler: userHandler,
		userService: userService,
		Db:          db,
	}
}
