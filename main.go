package main

import (
	"github.com/gin-gonic/gin"
	"watt-flow/config"
	"watt-flow/route"
	"watt-flow/server"
)

func main() {
	env := config.Init()
	dependencies := server.InitDeps(&env)
	gin.DefaultWriter = dependencies.Logger.GetGinLogger()
	engine := gin.New()

	route.RegisterRoutes(engine, dependencies)
	//middleware.RegisterMiddleware(engine, dependencies, env)
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	engine.Run(":" + env.ServerPort)

}
