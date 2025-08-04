package util

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// TransactionManager handles database transactions with proper connection lifecycle
type TransactionManager struct {
	db     *gorm.DB
	logger Logger
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *gorm.DB, logger Logger) *TransactionManager {
	return &TransactionManager{
		db:     db,
		logger: logger,
	}
}

// ExecuteInTransaction executes a function within a database transaction
// Automatically handles commit/rollback based on function success/failure
func (tm *TransactionManager) ExecuteInTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	// Begin transaction
	tx := tm.db.Begin()
	if tx.Error != nil {
		tm.logger.Error("Failed to begin transaction", tx.Error)
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Ensure transaction is properly closed
	defer func() {
		if r := recover(); r != nil {
			tm.logger.Error("Panic occurred in transaction, rolling back", fmt.Errorf("panic: %v", r))
			if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
				tm.logger.Error("Failed to rollback transaction after panic", rollbackErr)
			}
			panic(r) // Re-throw panic
		}
	}()

	// Execute the function
	if err := fn(tx); err != nil {
		tm.logger.Error("Transaction function failed, rolling back", err)
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			tm.logger.Error("Failed to rollback transaction", rollbackErr)
			return fmt.Errorf("transaction failed and rollback failed: %v, original error: %w", rollbackErr, err)
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tm.logger.Error("Failed to commit transaction", err)
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			tm.logger.Error("Failed to rollback after commit error", rollbackErr)
		}
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	tm.logger.Info("Transaction completed successfully")
	return nil
}

// GetTransaction returns a new transaction
func (tm *TransactionManager) GetTransaction() *gorm.DB {
	return tm.db.Begin()
}
