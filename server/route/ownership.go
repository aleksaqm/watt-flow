package route

import (
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/go-redis/redis/v8"
	"time"
	"watt-flow/middleware"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type OwnershipRoute struct {
	engine      *gin.Engine
	store       persist.CacheStore
	redisClient *redis.Client
}

func (r OwnershipRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up ownerships routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)

	getPendingOwnershipsRule := middleware.InvalidationRule{
		Pattern: "/api/ownership/pending*",
		Params:  []string{},
	}
	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.POST("/household/owner", authMid.RoleMiddleware([]string{"Regular"}), server.OwnershipHandler.CreateOwnershipRequest)
		api.GET("/ownership/requests/:id", authMid.RoleMiddleware([]string{"Regular"}), server.OwnershipHandler.GetOwnershipRequestsForUser)
		api.GET("/ownership/pending", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), cache.CacheByRequestURI(r.store, 2*time.Minute), server.OwnershipHandler.GetPendingRequests)
		api.PATCH("/ownership/accept/:id", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), middleware.CacheInvalidationMiddleware(r.store, r.redisClient, getPendingOwnershipsRule), server.OwnershipHandler.AcceptOwnershipRequest)
		api.PUT("/ownership/decline/:id", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), middleware.CacheInvalidationMiddleware(r.store, r.redisClient, getPendingOwnershipsRule), server.OwnershipHandler.DeclineOwnershipRequest)
	}
}

func NewOwnershipRoute(engine *gin.Engine, store persist.CacheStore, redisClient *redis.Client) *OwnershipRoute {
	return &OwnershipRoute{
		engine:      engine,
		store:       store,
		redisClient: redisClient,
	}
}
