package route

import (
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"time"
	"watt-flow/middleware"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type DeviceConsumptionRoute struct {
	engine *gin.Engine
	store  persist.CacheStore
}

func (r DeviceConsumptionRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up device consumption routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)

	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.POST("/device-consumption/query-consumption", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin", "Regular"}), cache.CacheByRequestURI(r.store, 2*time.Second), server.ElectricityConsumptionHandler.QueryConsumption)
	}
}

func NewDeviceConsumptionRoute(engine *gin.Engine, store persist.CacheStore) *DeviceConsumptionRoute {
	return &DeviceConsumptionRoute{
		engine: engine,
		store:  store,
	}
}
