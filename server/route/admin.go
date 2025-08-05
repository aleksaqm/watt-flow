package route

import (
	"github.com/gin-gonic/gin"
	"watt-flow/server"
)

type AdminRoute struct {
	engine *gin.Engine
}

func (r *AdminRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up admin routes")
	api := r.engine.Group("/api/admin")
	{
		api.POST("/password", server.UserHandler.ChangeAdminPassword)
		api.GET("/active", server.UserHandler.IsAdminActive)
	}
}

func NewAdminRoute(engine *gin.Engine) *AdminRoute {
	return &AdminRoute{
		engine: engine,
	}
}
