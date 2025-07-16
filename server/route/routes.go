package route

import (
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type Routes []Route

// Route interface
type Route interface {
	Register(server *server.Server)
}

func RegisterRoutes(
	engine *gin.Engine,
	server *server.Server,
) {
	NewUserRoute(engine).Register(server)
	NewAuthRoute(engine).Register(server)
	NewPropertyRoute(engine).Register(server)
	NewHouseholdRoute(engine).Register(server)
	NewOwnershipRoute(engine).Register(server)
	NewDeviceStatusRoute(engine).Register(server)
	NewAddressRoute(engine).Register(server)
	NewAdminRoute(engine).Register(server)
	NewPermissionRoute(engine).Register(server)
	NewCityRoute(engine).Register(server)
	NewMeetingRoute(engine).Register(server)
	NewPricelistRoute(engine).Register(server)
	NewBillRoute(engine).Register(server)
}
