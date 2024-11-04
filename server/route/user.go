package route

import (
	"github.com/gin-gonic/gin"
	"watt-flow/server"
)

type UserRoute struct {
	engine *gin.Engine
}

func (r UserRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up user routes")
	api := r.engine.Group("/api")
	{
		api.GET("/user", server.UserHandler.GetById)
	}
}

func NewUserRoute(engine *gin.Engine) *UserRoute {
	return &UserRoute{
		engine: engine,
	}
}
