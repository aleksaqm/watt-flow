package route

import (
	"github.com/gin-gonic/gin"
	"watt-flow/server"
)

type HouseholdRoute struct {
	engine *gin.Engine
}

func (r HouseholdRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up household routes")
	// authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)

	api := r.engine.Group("/api2")
	{
		api.GET("/household/:id", server.HouseholdHandler.GetById)
		api.POST("/household", server.HouseholdHandler.Create)
		api.PUT("/household/:id", server.HouseholdHandler.Update)
		api.GET("/households", server.HouseholdHandler.FindByStatus)
		api.DELETE("/household/:id", server.HouseholdHandler.Delete)
	}
}

func NewHouseholdRoute(engine *gin.Engine) *HouseholdRoute {
	return &HouseholdRoute{
		engine: engine,
	}
}
