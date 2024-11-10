package server

import (
	"watt-flow/db"
	"watt-flow/handler"
	"watt-flow/service"
	"watt-flow/util"
)

type Server struct {
	Logger           util.Logger
	UserHandler      *handler.UserHandler
	userService      service.IUserService
	PropertyHandler  *handler.PropertyHandler
	propertyService  service.IPropertyService
	HouseholdHandler *handler.HouseholdHandler
	householdService service.IHouseholdService
	AuthService      *service.AuthService
	Db               db.Database
}

func NewServer(logger util.Logger, userService service.IUserService, authService *service.AuthService, userHandler *handler.UserHandler,
	propertyService service.IPropertyService, propertyHandler *handler.PropertyHandler,
	householdService service.IHouseholdService, householdHandler *handler.HouseholdHandler, db db.Database) *Server {
	return &Server{
		Logger:           logger,
		UserHandler:      userHandler,
		userService:      userService,
		AuthService:      authService,
		PropertyHandler:  propertyHandler,
		propertyService:  propertyService,
		HouseholdHandler: householdHandler,
		householdService: householdService,
		Db:               db,
	}
}
