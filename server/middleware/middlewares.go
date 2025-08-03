package middleware

import (
	"github.com/gin-gonic/gin"
	"watt-flow/server"
)

type Middleware interface {
	Register()
}

type Middlewares []Middleware

func RegisterMiddlewares(
	engine *gin.Engine,
	server *server.Server,
) {

	// registering middlewares
	NewCorsMiddleware(engine, server.Logger).Register()
}
