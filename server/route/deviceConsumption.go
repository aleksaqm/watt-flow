package route

import (
	"watt-flow/middleware"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type DeviceConsumptionRoute struct {
	engine *gin.Engine
}

func (r DeviceConsumptionRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up device consumption routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)

	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.POST("/device-consumption/query-consumption", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin", "Regular"}), server.ElectricityConsumptionHandler.QueryConsumption)
	}
}

func NewDeviceConsumptionRoute(engine *gin.Engine) *DeviceConsumptionRoute {
	return &DeviceConsumptionRoute{
		engine: engine,
	}
}
