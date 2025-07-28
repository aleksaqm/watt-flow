package route

import (
	"watt-flow/middleware"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type HouseholdRoute struct {
	engine *gin.Engine
}

func (r HouseholdRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up household routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)

	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.GET("/household/:id", authMid.RoleMiddleware([]string{"Admin", "Clerk", "Regular", "SuperAdmin"}), server.HouseholdHandler.GetById)
		api.POST("/household/query", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin", "Regular"}), server.HouseholdHandler.Query)

		api.GET("/household/:id/consumption/monthly", authMid.RoleMiddleware([]string{"Regular"}), server.ElectricityConsumptionHandler.GetMonthlyConsumption)
		api.GET("/household/:id/consumption/12months", authMid.RoleMiddleware([]string{"Regular"}), server.ElectricityConsumptionHandler.Get12MonthsConsumption)
		api.GET("/household/:id/consumption/daily", authMid.RoleMiddleware([]string{"Regular"}), server.ElectricityConsumptionHandler.GetDailyConsumption)
	}
}

func NewHouseholdRoute(engine *gin.Engine) *HouseholdRoute {
	return &HouseholdRoute{
		engine: engine,
	}
}
