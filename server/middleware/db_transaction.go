package middleware

import (
	"net/http"
	"runtime/debug"
	"watt-flow/db"
	"watt-flow/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func (m DatabaseTrx) Register() {
	m.logger.Info("setting up database transaction middleware")

	m.engine.Use(func(c *gin.Context) {
		txHandle := m.db.DB.Begin()
		m.logger.Info("Transaction BEGIN")
		m.logger.Info("Transaction BEGIN", zap.String("path", c.Request.URL.Path))
		m.logger.Debug("Stack Trace:\n" + string(debug.Stack()))

		defer func() {
			if r := recover(); r != nil {
				txHandle.Rollback()
			}
		}()

		c.Set("db_trx", txHandle)
		c.Next()

		// commit transaction on success status
		if statusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated, http.StatusNoContent}) {
			m.logger.Info("Transaction COMMIT")
			if err := txHandle.Commit().Error; err != nil {
				m.logger.Error("trx commit error: ", err)
			}
		} else {
			m.logger.Info("rolling back transaction due to status code: 500")
			txHandle.Rollback()
		}
	})
}
