package repository

import (
	"watt-flow/db"
	"watt-flow/model"
	"watt-flow/util"
)

type OwnershipRepository struct {
	database db.Database
	logger   util.Logger
}

func NewOwnershipRepository(db db.Database, logger util.Logger) *OwnershipRepository {
	err := db.AutoMigrate(&model.OwnershipRequest{})
	if err != nil {
		logger.Error("Error migrating Ownership Repository", err)
	}
	return &OwnershipRepository{
		database: db,
		logger:   logger,
	}
}

func (repository *OwnershipRepository) Create(ownershipRequest *model.OwnershipRequest) (model.OwnershipRequest, error) {
	result := repository.database.Create(ownershipRequest)
	if result.Error != nil {
		repository.logger.Error("Error creating ownership request", result.Error)
		return *ownershipRequest, result.Error
	}
	return *ownershipRequest, nil
}
