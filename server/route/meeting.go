package route

import (
	"watt-flow/middleware"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type MeetingRoute struct {
	engine *gin.Engine
}

func (r MeetingRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up meeting routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)
	txMid := middleware.NewTransactionMiddleware(server.Logger, server.Db)

	api := r.engine.Group("/api").Use(authMid.Handler())
	{
		api.POST("/meeting", authMid.RoleMiddleware([]string{"Clerk", "Regular"}), txMid.Handler(), server.MeetingHandler.CreateNewMeeting)
		api.GET("/meeting/:id", authMid.RoleMiddleware([]string{"Clerk", "Regular"}), server.MeetingHandler.GetMeetingById)
		api.GET("/user/meetings/:id", authMid.RoleMiddleware([]string{"Clerk", "Regular"}), server.MeetingHandler.GetUsersMeetings)
		api.GET("/timeslot", authMid.RoleMiddleware([]string{"Clerk", "Regular"}), server.MeetingHandler.GetSlotByDateAndClerkId)
	}
}

func NewMeetingRoute(engine *gin.Engine) *MeetingRoute {
	return &MeetingRoute{
		engine: engine,
	}
}
