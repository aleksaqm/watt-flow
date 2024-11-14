package route

import (
	"github.com/gin-gonic/gin"
	"watt-flow/server"
)

type AddressRoute struct {
	engine *gin.Engine
}

func (r AddressRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up address routes")

	api := r.engine.Group("/api")
	{
		api.GET("/address/:id", server.AddressHandler.FindById)
		api.POST("/address", server.AddressHandler.Create)
		api.PUT("/address/:id", server.AddressHandler.Update)
		api.GET("/addresses", server.AddressHandler.FindAll)
		api.DELETE("/address/:id", server.AddressHandler.Delete)
	}
}

func NewAddressRoute(engine *gin.Engine) *AddressRoute {
	return &AddressRoute{
		engine: engine,
	}
}
