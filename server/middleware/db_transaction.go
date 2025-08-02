package middleware

import (
	"net/http"
	"strings"

	"watt-flow/db"
	"watt-flow/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DatabaseTrx struct {
	engine *gin.Engine
	logger util.Logger
	db     db.Database
}

func statusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}

func NewDatabaseTrx(
	engine *gin.Engine,
	logger util.Logger,
	db db.Database,
) DatabaseTrx {
	return DatabaseTrx{
		engine: engine,
		logger: logger,
		db:     db,
	}
}

// Check if the request needs a transaction
func needsTransaction(method, path string) bool {
	// Only transactional methods
	if method != "POST" && method != "PUT" && method != "DELETE" && method != "PATCH" {
		return false
	}

	// Skip health checks, metrics, etc.
	skipPaths := []string{"/health", "/metrics", "/ping", "/favicon.ico"}
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return false
		}
	}

	return true
}

func (m DatabaseTrx) Register() {
	m.logger.Info("setting up database transaction middleware")
	m.engine.Use(func(c *gin.Context) {
		// Check if this request actually needs a transaction
		if !needsTransaction(c.Request.Method, c.Request.URL.Path) {
			// For read-only requests, just provide the regular DB connection
			c.Set("db_trx", m.db.DB)
			c.Next()
			return
		}

		// Start transaction for write operations
		txHandle := m.db.DB.Begin()
		if txHandle.Error != nil {
			m.logger.Error("Failed to begin transaction", txHandle.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database transaction failed"})
			c.Abort()
			return
		}

		m.logger.Info("Transaction BEGIN", zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))

		// Ensure transaction is always closed
		var committed bool
		defer func() {
			if r := recover(); r != nil {
				if !committed {
					m.logger.Error("Transaction ROLLBACK (panic)", zap.String("path", c.Request.URL.Path), zap.Any("panic", r))
					txHandle.Rollback()
				}
				panic(r) // Re-panic after cleanup
			} else if !committed {
				// If we reach here without committing, something went wrong
				m.logger.Warn("Transaction ROLLBACK (not committed)", zap.String("path", c.Request.URL.Path), zap.Int("status", c.Writer.Status()))
				txHandle.Rollback()
			}
		}()

		c.Set("db_trx", txHandle)
		c.Next()

		// Commit or rollback based on result
		if c.IsAborted() || len(c.Errors) > 0 {
			m.logger.Info("Transaction ROLLBACK (aborted/errors)", zap.String("path", c.Request.URL.Path))
			txHandle.Rollback()
			committed = true
		} else if statusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated, http.StatusNoContent}) {
			m.logger.Info("Transaction COMMIT", zap.String("path", c.Request.URL.Path))
			if err := txHandle.Commit().Error; err != nil {
				m.logger.Error("Transaction commit error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
			}
			committed = true
		} else {
			m.logger.Info("Transaction ROLLBACK (status code)", zap.String("path", c.Request.URL.Path), zap.Int("status", c.Writer.Status()))
			txHandle.Rollback()
			committed = true
		}
	})
}

// Alternative: Selective transaction middleware
func (m DatabaseTrx) RegisterSelective() {
	m.logger.Info("setting up selective database transaction middleware")

	// Middleware that provides regular DB for most requests
	m.engine.Use(func(c *gin.Context) {
		c.Set("db", m.db.DB) // Always provide regular DB
		c.Next()
	})
}

// Helper function to manually manage transactions in handlers when needed
func WithTransaction(db *gorm.DB, logger util.Logger, fn func(*gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
