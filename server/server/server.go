package server

import (
	"watt-flow/db"
	"watt-flow/handler"
	"watt-flow/service"
	"watt-flow/util"
)

type Server struct {
	Logger              util.Logger
	UserHandler         *handler.UserHandler
	userService         service.IUserService
	PropertyHandler     *handler.PropertyHandler
	propertyService     service.IPropertyService
	HouseholdHandler    *handler.HouseholdHandler
	householdService    service.IHouseholdService
	DeviceStatusHandler *handler.DeviceStatusHandler
	deviceStatusService service.IDeviceStatusService
	AddressHandler      *handler.AddressHandler
	addressService      service.IAddressService
	CityHandler         *handler.CityHandler
	cityService         service.ICityService
	RestartService      *service.RestartService
	AuthService         *service.AuthService
	MeetingService      service.IMeetingService
	MeetingHandler      *handler.MeetingHandler
	Db                  db.Database
}

func NewServer(logger util.Logger, userService service.IUserService, authService *service.AuthService, restartService *service.RestartService, userHandler *handler.UserHandler,
	propertyService service.IPropertyService, propertyHandler *handler.PropertyHandler,
	householdService service.IHouseholdService, householdHandler *handler.HouseholdHandler,
	deviceStatusService service.IDeviceStatusService, deviceStatusHandler *handler.DeviceStatusHandler,
	addressService service.IAddressService, addressHandler *handler.AddressHandler,
	meetingService service.IMeetingService, meetingHandler *handler.MeetingHandler,
	cityService service.ICityService, cityHandler *handler.CityHandler, db db.Database) *Server {
	return &Server{
		Logger:              logger,
		UserHandler:         userHandler,
		RestartService:      restartService,
		userService:         userService,
		AuthService:         authService,
		PropertyHandler:     propertyHandler,
		propertyService:     propertyService,
		HouseholdHandler:    householdHandler,
		householdService:    householdService,
		DeviceStatusHandler: deviceStatusHandler,
		deviceStatusService: deviceStatusService,
		addressService:      addressService,
		AddressHandler:      addressHandler,
		cityService:         cityService,
		CityHandler:         cityHandler,
		MeetingService:      meetingService,
		MeetingHandler:      meetingHandler,
		Db:                  db,
	}
}
