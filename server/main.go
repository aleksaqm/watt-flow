package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"watt-flow/config"
	"watt-flow/middleware"
	"watt-flow/route"
	"watt-flow/server"
	"watt-flow/util"

	"github.com/gin-gonic/gin"
)

func main() {
	env := config.Init()
	dependencies := server.InitDeps(env)
	gin.DefaultWriter = dependencies.Logger.GetGinLogger()
	engine := gin.New()

	engine.Use(gin.Recovery())
	// engine.Use(gin.Logger())

	middleware.RegisterMiddlewares(engine, dependencies)
	route.RegisterRoutes(engine, dependencies)

	if env.Restart {
		err := dependencies.RestartService.ResetDatabase()
		if err != nil {
			dependencies.Logger.Error("Error resetting database", err)
		}
		err = dependencies.RestartService.InitSuperAdmin()
		if err != nil {
			dependencies.Logger.Error("Error initializing super admin", err)
		} else {
			dependencies.Logger.Info("Database reset and super admin initialized successfully")
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go monitorDBConnections(ctx, dependencies, dependencies.Logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		dependencies.Logger.Info("Shutting down server...")
		cancel()
		os.Exit(0)
	}()

	engine.Run(":" + env.ServerPort)
}

func monitorDBConnections(ctx context.Context, dependencies *server.Server, logger util.Logger) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping DB connection monitoring")
			return
		case <-ticker.C:
			// Get the underlying sql.DB instance
			sqlDB, err := dependencies.Db.DB.DB()
			if err != nil {
				logger.Error("Failed to get database instance for monitoring", err)
				continue
			}

			stats := sqlDB.Stats()
			logger.Info("DB Connection Stats", map[string]interface{}{
				"open_connections": stats.OpenConnections,
				"in_use":           stats.InUse,
				"idle":             stats.Idle,
				"wait_count":       stats.WaitCount,
				"wait_duration":    stats.WaitDuration.String(),
				"max_idle":         stats.MaxIdleClosed,
				"max_lifetime":     stats.MaxLifetimeClosed,
			})

			// Alert if connection usage is high
			if stats.InUse > 20 { // Adjust threshold as needed
				logger.Warn("High database connection usage detected", map[string]interface{}{
					"in_use":           stats.InUse,
					"open_connections": stats.OpenConnections,
				})
			}

			// Alert if there are connection waits
			if stats.WaitCount > 0 {
				logger.Warn("Database connection waits detected", map[string]interface{}{
					"wait_count":    stats.WaitCount,
					"wait_duration": stats.WaitDuration.String(),
				})
			}
		}
	}
}
