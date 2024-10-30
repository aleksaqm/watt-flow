// go:build wireinject
//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"
	"watt-flow/config"
	"watt-flow/db"
	"watt-flow/handler"
	"watt-flow/repository"
	"watt-flow/service"
	"watt-flow/util"
)

var userServiceSet = wire.NewSet(service.NewUserService, wire.Bind(new(service.IUserService), new(*service.UserService)))

func InitDeps(env *config.Environment) *Server {
	wire.Build(db.NewDatabase, util.NewLogger, repository.NewUserRepository, userServiceSet, handler.NewUserHandler, NewServer)
	return &Server{}
}
