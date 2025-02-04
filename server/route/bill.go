package route

import (
	"watt-flow/middleware"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type BillRoute struct {
	engine *gin.Engine
}

func (r BillRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up bill routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)
	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.GET("/bills/unsent", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.BillHandler.GetUnsentMonthlyBills)
		api.POST("/bills/send", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.BillHandler.InitiateBilling)
		api.GET("/bills/query", authMid.RoleMiddleware([]string{"Admin", "SuperAdmin"}), server.BillHandler.Query)
	}
}

func NewBillRoute(engine *gin.Engine) *BillRoute {
	return &BillRoute{
		engine: engine,
	}
}
