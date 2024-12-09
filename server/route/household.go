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
		api.POST("/household/owner", authMid.RoleMiddleware([]string{"Regular"}), server.HouseholdHandler.CreateOwnershipRequest)
		api.GET("/ownership/requests/:id", authMid.RoleMiddleware([]string{"Regular"}), server.HouseholdHandler.GetOwnershipRequestsForUser)
		api.GET("/ownership/pending", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.HouseholdHandler.GetPendingRequests)
		api.PATCH("/ownership/accept/:id", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.HouseholdHandler.AcceptOwnershipRequest)
	}
}

func NewHouseholdRoute(engine *gin.Engine) *HouseholdRoute {
	return &HouseholdRoute{
		engine: engine,
	}
}
