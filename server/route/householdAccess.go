package route

import (
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"time"
	"watt-flow/middleware"
	"watt-flow/server"
)

type HouseholdAccessRoute struct {
	engine *gin.Engine
	store  persist.CacheStore
}

func (r HouseholdAccessRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up household access routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)

	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.POST("/household/:id/access", authMid.RoleMiddleware([]string{"Admin", "Clerk", "Regular", "SuperAdmin"}), server.HouseholdAccessHandler.GrantAccess)
		api.DELETE("/household/:householdId/access/revoke/:userId", authMid.RoleMiddleware([]string{"Admin", "Clerk", "Regular", "SuperAdmin"}), server.HouseholdAccessHandler.RevokeAccess)
		api.GET("/household/access/:householdId", authMid.RoleMiddleware([]string{"Admin", "Clerk", "Regular", "SuperAdmin"}), cache.CacheByRequestURI(r.store, 2*time.Second), server.HouseholdAccessHandler.ListAccess)
	}
}

func NewHouseholdAccessRoute(engine *gin.Engine, store persist.CacheStore) *HouseholdAccessRoute {
	return &HouseholdAccessRoute{
		engine: engine,
		store:  store,
	}
}
