package route

import (
	"github.com/gin-gonic/gin"
	"watt-flow/server"
)

type PropertyRoute struct {
	engine *gin.Engine
}

func (r PropertyRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up property routes")
	//authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)

	api := r.engine.Group("/api")
	{
		api.GET("/property/:id", server.PropertyHandler.GetById)
		api.POST("/property", server.PropertyHandler.Create)
		api.PUT("/property/:id", server.PropertyHandler.Update)
		api.GET("/properties", server.PropertyHandler.FindByStatus)
		api.DELETE("/property/:id", server.PropertyHandler.Delete)
	}
}

func NewPropertyRoute(engine *gin.Engine) *PropertyRoute {
	return &PropertyRoute{
		engine: engine,
	}
}
