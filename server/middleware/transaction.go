package middleware

import (
	"net/http"

	"watt-flow/db"
	"watt-flow/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TransactionMiddleware struct {
	logger util.Logger
	db     db.Database
}

func NewTransactionMiddleware(logger util.Logger, db db.Database) *TransactionMiddleware {
	return &TransactionMiddleware{
		logger: logger,
		db:     db,
	}
}

func (m *TransactionMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		txHandle := m.db.DB.Begin()
		if txHandle.Error != nil {
			m.logger.Error("Failed to begin transaction", txHandle.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			c.Abort()
			return
		}

		m.logger.Info("Transaction BEGIN", zap.String("path", c.Request.URL.Path))

		// Ensure transaction is always closed
		var txClosed bool
		defer func() {
			if r := recover(); r != nil {
				if !txClosed {
					m.logger.Error("Panic occurred, rolling back transaction", zap.Any("panic", r))
					if err := txHandle.Rollback().Error; err != nil {
						m.logger.Error("Failed to rollback transaction after panic", err)
					}
					txClosed = true
				}
				panic(r) // Re-throw panic
			}

			// Final safety check - if transaction somehow wasn't closed
			if !txClosed {
				m.logger.Warn("Transaction was not properly closed, forcing rollback",
					zap.String("path", c.Request.URL.Path))
				if err := txHandle.Rollback().Error; err != nil {
					m.logger.Error("Failed to rollback transaction in final safety check", err)
				}
			}
		}()

		c.Set("db_trx", txHandle)
		c.Next()

		// Check if transaction was already handled
		if txClosed {
			return
		}

		status := c.Writer.Status()

		// Commit on success status codes
		if statusInList(status, []int{http.StatusOK, http.StatusCreated, http.StatusNoContent}) {
			m.logger.Info("Transaction COMMIT", zap.String("path", c.Request.URL.Path), zap.Int("status", status))
			if err := txHandle.Commit().Error; err != nil {
				m.logger.Error("Transaction commit error", err)
				// If commit fails, try to rollback
				if rollbackErr := txHandle.Rollback().Error; rollbackErr != nil {
					m.logger.Error("Failed to rollback after commit error", rollbackErr)
				}
			}
			txClosed = true
		} else {
			// Rollback on any error status code (400, 401, 403, 404, 422, 500, etc.)
			m.logger.Info("Transaction ROLLBACK", zap.String("path", c.Request.URL.Path), zap.Int("status", status))
			if err := txHandle.Rollback().Error; err != nil {
				m.logger.Error("Transaction rollback error", err)
			}
			txClosed = true
		}
	}
}

func statusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}