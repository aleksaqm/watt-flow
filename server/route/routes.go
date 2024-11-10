package route

import (
	"github.com/gin-gonic/gin"
	"watt-flow/server"
)

type Routes []Route

// Route interface
type Route interface {
	Register(server *server.Server)
}

func RegisterRoutes(
	engine *gin.Engine,
	server *server.Server,
) {
	NewUserRoute(engine).Register(server)
	NewAuthRoute(engine).Register(server)
	NewPropertyRoute(engine).Register(server)
	NewHouseholdRoute(engine).Register(server)
}
