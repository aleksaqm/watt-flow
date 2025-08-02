package route

import (
	"github.com/chenyahui/gin-cache/persist"
	"watt-flow/middleware"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type OwnershipRoute struct {
	engine *gin.Engine
	store  persist.CacheStore
}

func (r OwnershipRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up ownerships routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)

	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.POST("/household/owner", authMid.RoleMiddleware([]string{"Regular"}), server.OwnershipHandler.CreateOwnershipRequest)
		api.GET("/ownership/requests/:id", authMid.RoleMiddleware([]string{"Regular"}), server.OwnershipHandler.GetOwnershipRequestsForUser)
		api.GET("/ownership/pending", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.OwnershipHandler.GetPendingRequests)
		api.PATCH("/ownership/accept/:id", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.OwnershipHandler.AcceptOwnershipRequest)
		api.PUT("/ownership/decline/:id", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.OwnershipHandler.DeclineOwnershipRequest)
	}
}

func NewOwnershipRoute(engine *gin.Engine, store persist.CacheStore) *OwnershipRoute {
	return &OwnershipRoute{
		engine: engine,
		store:  store,
	}
}
