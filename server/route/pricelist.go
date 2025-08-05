package route

import (
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"time"
	"watt-flow/middleware"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type PricelistRoute struct {
	engine *gin.Engine
	store  persist.CacheStore
}

func (r PricelistRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up pricelist routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)
	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.POST("/pricelist", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.PricelistHandler.CreatePricelist)
		api.GET("/pricelist/query", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), cache.CacheByRequestURI(r.store, 2*time.Second), server.PricelistHandler.Query)
		api.DELETE("/pricelist/:id", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.PricelistHandler.Delete)
	}
}

func NewPricelistRoute(engine *gin.Engine, store persist.CacheStore) *PricelistRoute {
	return &PricelistRoute{
		engine: engine,
		store:  store,
	}
}
