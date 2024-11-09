package main

import (
	"github.com/gin-gonic/gin"
	"watt-flow/config"
	"watt-flow/middleware"
	"watt-flow/route"
	"watt-flow/server"
)

func main() {
	env := config.Init()
	dependencies := server.InitDeps(env)
	gin.DefaultWriter = dependencies.Logger.GetGinLogger()
	engine := gin.New()

	route.RegisterRoutes(engine, dependencies)
	middleware.RegisterMiddlewares(engine, dependencies)
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	if env.Restart {
		_ = dependencies.RestartService.ResetDatabase()
		_ = dependencies.RestartService.InitSuperAdmin()
	}

	//dependencies.UserService.InitSuperAdmin()
	engine.Run(":" + env.ServerPort)
}
