package route

import (
	"watt-flow/middleware"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type PricelistRoute struct {
	engine *gin.Engine
}

func (r PricelistRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up pricelist routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)
	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.POST("/pricelist", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.PricelistHandler.CreatePricelist)
		api.GET("/pricelist/query", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.PricelistHandler.Query)
		api.DELETE("/pricelist/:id", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.PricelistHandler.Delete)
	}
}

func NewPricelistRoute(engine *gin.Engine) *PricelistRoute {
	return &PricelistRoute{
		engine: engine,
	}
}
