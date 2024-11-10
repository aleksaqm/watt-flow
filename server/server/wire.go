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

var propertyServiceSet = wire.NewSet(
	service.NewPropertyService,
	wire.Bind(new(service.IPropertyService), new(*service.PropertyService)),
)

var householdServiceSet = wire.NewSet(
	service.NewHouseholdService,
	wire.Bind(new(service.IHouseholdService), new(*service.HouseholdService)))

var deviceStatusServiceSet = wire.NewSet(
	service.NewDeviceStatusService,
	wire.Bind(new(service.IDeviceStatusService), new(*service.DeviceStatusService)))

func InitDeps(env *config.Environment) *Server {
	wire.Build(db.NewDatabase, util.NewLogger, repository.NewUserRepository, service.NewAuthService, userServiceSet, handler.NewUserHandler,
		repository.NewPropertyRepository, propertyServiceSet, handler.NewPropertyHandler,
		repository.NewHouseholdRepository, householdServiceSet, handler.NewHouseholdHandler,
		repository.NewDeviceStatusRepository, deviceStatusServiceSet, handler.NewDeviceStatusHandler, NewServer)
	return &Server{}
}
