package route

import (
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"time"
	"watt-flow/middleware"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type MeetingRoute struct {
	engine *gin.Engine
	store  persist.CacheStore
}

func (r MeetingRoute) Register(server *server.Server) {
	server.Logger.Info("Setting up meeting routes")
	authMid := middleware.NewAuthMiddleware(server.AuthService, server.Logger)

	api := r.engine.Group("/api")
	{
		api.POST("/meeting", authMid.RoleMiddleware([]string{"Clerk", "Regular"}), server.MeetingHandler.CreateNewMeeting)
		api.GET("/meeting/:id", authMid.RoleMiddleware([]string{"Clerk", "Regular"}), cache.CacheByRequestURI(r.store, 2*time.Second), server.MeetingHandler.GetMeetingById)
		api.GET("/user/meetings/:id", authMid.RoleMiddleware([]string{"Clerk", "Regular"}), cache.CacheByRequestURI(r.store, 2*time.Second), server.MeetingHandler.GetUsersMeetings)
		api.GET("/timeslot", authMid.RoleMiddleware([]string{"Clerk", "Regular"}), cache.CacheByRequestURI(r.store, 2*time.Second), server.MeetingHandler.GetSlotByDateAndClerkId)
	}
}

func NewMeetingRoute(engine *gin.Engine, store persist.CacheStore) *MeetingRoute {
	return &MeetingRoute{
		engine: engine,
		store:  store,
	}
}
