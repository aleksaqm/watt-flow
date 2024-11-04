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

	route.RegisterRoutes(engine, dependencies)
	middleware.RegisterMiddlewares(engine, dependencies)
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	engine.Run(":" + env.ServerPort)
}
