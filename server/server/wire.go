// go:build wireinject
//go:build wireinject
// +build wireinject

package server

import (
	"watt-flow/config"
	"watt-flow/db"
	"watt-flow/handler"
	"watt-flow/repository"
	"watt-flow/service"
	"watt-flow/util"

	"github.com/google/wire"
)

var userServiceSet = wire.NewSet(service.NewUserService, wire.Bind(new(service.IUserService), new(*service.UserService)))

var propertyServiceSet = wire.NewSet(
	service.NewPropertyService,
	wire.Bind(new(service.IPropertyService), new(*service.PropertyService)),
)

var householdServiceSet = wire.NewSet(
	service.NewHouseholdService,
	wire.Bind(new(service.IHouseholdService), new(*service.HouseholdService)))

var ownershipServiceSet = wire.NewSet(
	service.NewOwnershipService,
	wire.Bind(new(service.IOwnershipService), new(*service.OwnershipService)))

var deviceStatusServiceSet = wire.NewSet(
	service.NewDeviceStatusService,
	wire.Bind(new(service.IDeviceStatusService), new(*service.DeviceStatusService)))

var addressServiceSet = wire.NewSet(
	service.NewAddressService,
	wire.Bind(new(service.IAddressService), new(*service.AddressService)))

var cityServiceSet = wire.NewSet(
	service.NewCityService,
	wire.Bind(new(service.ICityService), new(*service.CityService)))

var meetingServiceSet = wire.NewSet(
	service.NewMeetingService,
	wire.Bind(new(service.IMeetingService), new(*service.MeetingService)))

var pricelistServiceSet = wire.NewSet(
	service.NewPricelistService,
	wire.Bind(new(service.IPricelistService), new(*service.PricelistService)))

var billServiceSet = wire.NewSet(
	service.NewBillService,
	wire.Bind(new(service.IBillService), new(*service.BillService)))

var electricityConsumptionServiceSet = wire.NewSet(
	service.NewElectricityConsumptionService)

func InitDeps(env *config.Environment) *Server {
	wire.Build(db.NewDatabase, util.NewLogger, util.NewEmailSender, util.NewInfluxQueryHelper, repository.NewUserRepository, service.NewAuthService, userServiceSet, service.NewRestartService, handler.NewUserHandler,
		repository.NewPropertyRepository, propertyServiceSet, handler.NewPropertyHandler,
		repository.NewHouseholdRepository, householdServiceSet, handler.NewHouseholdHandler,
		repository.NewOwnershipRepository, ownershipServiceSet, handler.NewOwnershipHandler,
		repository.NewPricelistRepository, pricelistServiceSet, handler.NewPricelistHandler,
		repository.NewBillRepository, repository.NewMonthlyBillRepository, billServiceSet, handler.NewBillHandler,
		repository.NewDeviceStatusRepository, deviceStatusServiceSet, handler.NewDeviceStatusHandler,
		repository.NewAddressRepository, addressServiceSet, handler.NewAddressHandler,
		repository.NewCityRepository, cityServiceSet, handler.NewCityHandler,
		repository.NewTimeSlotRepository, repository.NewMeetingRepository, meetingServiceSet, handler.NewMeetingHandler,
		electricityConsumptionServiceSet, handler.NewElectricityConsumptionHandler, NewServer)
	return &Server{}
}
