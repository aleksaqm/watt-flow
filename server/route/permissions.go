package route

import (
	"github.com/gin-gonic/gin"
	"watt-flow/middleware"
	"watt-flow/server"
)

type PermissionRoute struct {
	engine *gin.Engine
}

func (r PermissionRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up user routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)
	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.GET("/validate/admin", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.UserHandler.ReturnOk)
	}
}

func NewPermissionRoute(engine *gin.Engine) *PermissionRoute {
	return &PermissionRoute{
		engine: engine,
	}
}
