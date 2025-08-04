package route

import (
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"time"
	"watt-flow/server"
)

type CityRoute struct {
	engine *gin.Engine
	store  persist.CacheStore
}

func (r CityRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up address city routes")

	api := r.engine.Group("/api")
	{
		api.GET("/cities", cache.CacheByRequestURI(r.store, 5*time.Minute), server.CityHandler.GetAllCities)
	}
}

func NewCityRoute(engine *gin.Engine, store persist.CacheStore) *CityRoute {
	return &CityRoute{
		engine: engine,
		store:  store,
	}
}
