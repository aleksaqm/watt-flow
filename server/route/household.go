package route

import (
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"time"
	"watt-flow/middleware"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type HouseholdRoute struct {
	engine *gin.Engine
	store  persist.CacheStore
}

func (r HouseholdRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up household routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)

	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.GET("/household/:id", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin", "Regular"}), cache.CacheByRequestURI(r.store, 2*time.Minute), server.HouseholdHandler.GetById)
		api.POST("/household/query", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin", "Regular"}), server.HouseholdHandler.Query)

		api.GET("/household/:id/consumption/monthly", authMid.RoleMiddleware([]string{"Regular"}), cache.CacheByRequestURI(r.store, 2*time.Second), server.ElectricityConsumptionHandler.GetMonthlyConsumption)
		api.GET("/household/:id/consumption/12months", authMid.RoleMiddleware([]string{"Regular"}), cache.CacheByRequestURI(r.store, 2*time.Second), server.ElectricityConsumptionHandler.Get12MonthsConsumption)
		api.GET("/household/:id/consumption/daily", authMid.RoleMiddleware([]string{"Regular"}), cache.CacheByRequestURI(r.store, 2*time.Second), server.ElectricityConsumptionHandler.GetDailyConsumption)

	}
}

func NewHouseholdRoute(engine *gin.Engine, store persist.CacheStore) *HouseholdRoute {
	return &HouseholdRoute{
		engine: engine,
		store:  store,
	}
}
