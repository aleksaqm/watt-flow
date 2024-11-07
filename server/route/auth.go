package route

import (
	"github.com/gin-gonic/gin"
	"watt-flow/server"
)

type AuthRoute struct {
	engine *gin.Engine
}

func (r *AuthRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up auth routes")
	api := r.engine.Group("/")
	{
		api.POST("login", server.UserHandler.Login)
		api.POST("register", server.UserHandler.Register)
	}
}

func NewAuthRoute(engine *gin.Engine) *AuthRoute {
	return &AuthRoute{engine: engine}
}
