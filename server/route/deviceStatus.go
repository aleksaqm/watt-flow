package route

import (
	"watt-flow/middleware"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type DeviceStatusRoute struct {
	engine *gin.Engine
}

func (r DeviceStatusRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up device status routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)

	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		// api.GET("/device-status/:address", server.DeviceStatusHandler.GetByAddress)
		// api.GET("/device-status/household/:household_id", server.DeviceStatusHandler.GetByHouseholdID)
		api.POST("/device-status/query-status", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.DeviceStatusHandler.QueryStatus)
		api.POST("/device-status/query-consumption", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.DeviceStatusHandler.QueryConsumption)
	}
}

func NewDeviceStatusRoute(engine *gin.Engine) *DeviceStatusRoute {
	return &DeviceStatusRoute{
		engine: engine,
	}
}
