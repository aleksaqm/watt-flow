package route

import (
	"watt-flow/middleware"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	engine *gin.Engine
}

func (r UserRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up user routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)
	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.GET("/user/:id", server.UserHandler.GetById)
		api.POST("/user", server.UserHandler.Create)
		api.POST("/user/clerk/new", authMid.RoleMiddleware([]string{"SuperAdmin", "Admin"}), server.UserHandler.RegisterClerk)
		api.GET("/user/admins", authMid.RoleMiddleware([]string{"SuperAdmin"}), server.UserHandler.FindAdmins)
		api.POST("/user/admin", authMid.RoleMiddleware([]string{"SuperAdmin"}), server.UserHandler.Create)
	}
}

func NewUserRoute(engine *gin.Engine) *UserRoute {
	return &UserRoute{
		engine: engine,
	}
}
