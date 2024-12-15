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

	api := r.engine.Group("/api")
	{
		api.PUT("/timeslot", authMid.RoleMiddleware([]string{"Clerk", "Regular"}), server.MeetingHandler.CreateOrUpdateTimeSlot)
		api.GET("/timeslot", authMid.RoleMiddleware([]string{"Clerk", "Regular"}), server.MeetingHandler.GetSlotById)
	}
}

func NewMeetingRoute(engine *gin.Engine) *MeetingRoute {
	return &MeetingRoute{
		engine: engine,
	}
}
