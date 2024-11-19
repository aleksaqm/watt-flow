package route

import (
	"github.com/gin-gonic/gin"
	"watt-flow/server"
)

type CityRoute struct {
	engine *gin.Engine
}

func (r CityRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up address city routes")

	api := r.engine.Group("/api")
	{
		api.GET("/cities", server.CityHandler.GetAllCities)
	}
}

func NewCityRoute(engine *gin.Engine) *CityRoute {
	return &CityRoute{
		engine: engine,
	}
}
