package route

import (
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"time"
	"watt-flow/middleware"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type BillRoute struct {
	engine *gin.Engine
	store  persist.CacheStore
}

func (r BillRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up bill routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)
	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.GET("/bills/unsent", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), cache.CacheByRequestURI(r.store, 2*time.Second), server.BillHandler.GetUnsentMonthlyBills)
		api.POST("/bills/send", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.BillHandler.InitiateBilling)
		api.GET("/bills/query", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), cache.CacheByRequestURI(r.store, 2*time.Second), server.BillHandler.Query)
		api.GET("/bills/search", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin", "Regular"}), cache.CacheByRequestURI(r.store, 2*time.Second), server.BillHandler.SearchBills)
		api.GET("/bills", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin", "Regular"}), cache.CacheByRequestURI(r.store, 2*time.Second), server.BillHandler.GetBill)
		api.PUT("/bills/pay/:id", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin", "Regular"}), server.BillHandler.PayBill)
	}
}

func NewBillRoute(engine *gin.Engine, store persist.CacheStore) *BillRoute {
	return &BillRoute{
		engine: engine,
		store:  store,
	}
}
