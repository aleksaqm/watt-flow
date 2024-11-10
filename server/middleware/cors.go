package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"watt-flow/util"
)

type CorsMiddleware struct {
	engine *gin.Engine
	logger util.Logger
}

func NewCorsMiddleware(engine *gin.Engine, logger util.Logger) CorsMiddleware {
	return CorsMiddleware{
		engine: engine,
		logger: logger,
	}
}

func (m CorsMiddleware) Register() {
	m.logger.Info("Setting up CORS middleware")
	m.engine.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:5173"},
			AllowCredentials: true,
			//AllowOriginFunc:  func(origin string) bool { return true }, // here origins need to be configured
			AllowHeaders: []string{"*"},
			//AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowMethods: []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"},
		}))
}
