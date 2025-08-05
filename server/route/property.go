package route

import (
	"time"
	"watt-flow/middleware"
	"watt-flow/server"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"

	"github.com/gin-gonic/gin"
)

type PropertyRoute struct {
	engine *gin.Engine
	store  persist.CacheStore
}

func (r PropertyRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up property routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)
	txMid := middleware.NewTransactionMiddleware(server.Logger, server.Db)

	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		// api.GET("/property/:id", server.PropertyHandler.GetById)
		api.POST("/property", authMid.RoleMiddleware([]string{"Regular"}), server.PropertyHandler.Create)
		// api.PUT("/property/:id", server.PropertyHandler.Update)
		// api.GET("/properties", server.PropertyHandler.FindByStatus)
		// api.DELETE("/property/:id", server.PropertyHandler.Delete)
		api.GET("/property/query", authMid.RoleMiddleware([]string{"Regular", "Admin", "SuperAdmin", "Clerk"}), cache.CacheByRequestURI(r.store, 2*time.Second), server.PropertyHandler.TableQuery)
		api.PUT("/property/:id/accept", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), txMid.Handler(), server.PropertyHandler.AcceptProperty)
		api.PUT("/property/:id/decline", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.PropertyHandler.DeclineProperty)
	}
}

func NewPropertyRoute(engine *gin.Engine, store persist.CacheStore) *PropertyRoute {
	return &PropertyRoute{
		engine: engine,
		store:  store,
	}
}
