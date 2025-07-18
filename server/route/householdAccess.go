package route

import (
	"github.com/gin-gonic/gin"
	"watt-flow/middleware"
	"watt-flow/server"
)

type HouseholdAccessRoute struct {
	engine *gin.Engine
}

func (r HouseholdAccessRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up household access routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)

	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.POST("/household/:id/access", authMid.RoleMiddleware([]string{"Admin", "Clerk", "Regular", "SuperAdmin"}), server.HouseholdAccessHandler.GrantAccess)
		api.DELETE("/household/:householdId/access/revoke/:userId", authMid.RoleMiddleware([]string{"Admin", "Clerk", "Regular", "SuperAdmin"}), server.HouseholdAccessHandler.RevokeAccess)

	}
}

func NewHouseholdAccessRoute(engine *gin.Engine) *HouseholdAccessRoute {
	return &HouseholdAccessRoute{
		engine: engine,
	}
}
