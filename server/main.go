package main

import (
	"watt-flow/config"
	"watt-flow/middleware"
	"watt-flow/route"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

func main() {
	env := config.Init()
	dependencies := server.InitDeps(env)
	gin.DefaultWriter = dependencies.Logger.GetGinLogger()
	engine := gin.New()

	middleware.RegisterMiddlewares(engine, dependencies)
	route.RegisterRoutes(engine, dependencies)
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	if env.Restart {
		err := dependencies.RestartService.ResetDatabase()
		if err != nil {
			dependencies.Logger.Error("Error resetting database", err)
		}
		err = dependencies.RestartService.InitSuperAdmin()
		if err != nil {
			dependencies.Logger.Error("Error initializing super admin", err)
		} else {
			dependencies.Logger.Info("Database reset and super admin initialized successfully")
		}
	}

	engine.Run(":" + env.ServerPort)
}
