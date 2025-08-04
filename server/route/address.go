package route

import (
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"time"
	"watt-flow/server"
)

type AddressRoute struct {
	engine *gin.Engine
	store  persist.CacheStore
}

func (r AddressRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up address routes")

	api := r.engine.Group("/api")
	{
		api.GET("/address/:id", cache.CacheByRequestURI(r.store, 10*time.Second), server.AddressHandler.FindById)
		api.POST("/address", server.AddressHandler.Create)
		api.PUT("/address/:id", server.AddressHandler.Update)
		api.GET("/addresses", cache.CacheByRequestURI(r.store, 10*time.Second), server.AddressHandler.FindAll)
		api.DELETE("/address/:id", server.AddressHandler.Delete)
	}
}

func NewAddressRoute(engine *gin.Engine, store persist.CacheStore) *AddressRoute {
	return &AddressRoute{
		engine: engine,
		store:  store,
	}
}
