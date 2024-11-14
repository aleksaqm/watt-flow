package route

import (
	"github.com/gin-gonic/gin"
	"watt-flow/server"
)

type DeviceStatusRoute struct {
	engine *gin.Engine
}

func (r DeviceStatusRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up device status routes")

	api := r.engine.Group("/api")
	{
		api.GET("/device-status/:address", server.DeviceStatusHandler.GetByAddress)
		api.GET("/device-status/household/:household_id", server.DeviceStatusHandler.GetByHouseholdID)
		api.POST("/device-status", server.DeviceStatusHandler.Create)
		api.PUT("/device-status/:address", server.DeviceStatusHandler.Update)
		api.DELETE("/device-status/:address", server.DeviceStatusHandler.Delete)
	}
}

func NewDeviceStatusRoute(engine *gin.Engine) *DeviceStatusRoute {
	return &DeviceStatusRoute{
		engine: engine,
	}
}
