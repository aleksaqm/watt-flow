package route

import (
	"context"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"watt-flow/server"

	"github.com/gin-gonic/gin"
)

type Routes []Route

type Route interface {
	Register(server *server.Server)
}

func RegisterRoutes(
	engine *gin.Engine,
	server *server.Server,
) {
	log.Println("REDDIS ADDR: ", os.Getenv("REDIS_ADDR"))
	log.Println("REDDIS PASSWORD: ", os.Getenv("REDIS_PASSWORD"))
	redisClientForCache := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       1,
	})

	if _, err := redisClientForCache.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("FATAL: Could not connect to Redis for caching: %v", err)
	}
	log.Println("Successfully connected to Redis DB 1 for caching.")

	cacheStore := persist.NewRedisStore(redisClientForCache)

	NewUserRoute(engine, cacheStore).Register(server)
	NewAuthRoute(engine).Register(server)
	NewPropertyRoute(engine, cacheStore).Register(server)
	NewHouseholdRoute(engine, cacheStore).Register(server)
	NewOwnershipRoute(engine, cacheStore, redisClientForCache).Register(server)
	NewDeviceStatusRoute(engine, cacheStore).Register(server)
	NewDeviceConsumptionRoute(engine, cacheStore).Register(server)
	NewAddressRoute(engine, cacheStore).Register(server)
	NewAdminRoute(engine).Register(server)
	NewPermissionRoute(engine).Register(server)
	NewCityRoute(engine, cacheStore).Register(server)
	NewMeetingRoute(engine, cacheStore).Register(server)
	NewPricelistRoute(engine, cacheStore).Register(server)
	NewBillRoute(engine, cacheStore).Register(server)
	NewHouseholdAccessRoute(engine, cacheStore).Register(server)
}
