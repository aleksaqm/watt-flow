package repository

import (
	"fmt"
	"gorm.io/gorm"
	"time"
	"watt-flow/db"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/util"
)

type OwnershipRepository struct {
	Database db.Database
	Logger   util.Logger
}

func NewOwnershipRepository(db db.Database, logger util.Logger) *OwnershipRepository {
	err := db.AutoMigrate(&model.OwnershipRequest{})
	if err != nil {
		logger.Error("Error migrating Ownership Repository", err)
	}
	return &OwnershipRepository{
		Database: db,
		Logger:   logger,
	}
}

func (repository *OwnershipRepository) Create(ownershipRequest *model.OwnershipRequest) (model.OwnershipRequest, error) {
	result := repository.Database.Create(ownershipRequest)
	if result.Error != nil {
		repository.Logger.Error("Error creating ownership request", result.Error)
		return *ownershipRequest, result.Error
	}
	return *ownershipRequest, nil
}

func (repository *OwnershipRepository) FindForOwner(ownerId uint64, params *dto.OwnershipQueryParams) ([]model.OwnershipRequest, int64, error) {
	var ownershipRequests []model.OwnershipRequest
	var total int64
	baseQuery := repository.Database.Model(&model.OwnershipRequest{}).
		Joins("JOIN households ON households.id = ownership_requests.household_id").
		Joins("JOIN properties ON properties.id = households.property_id").
		Where("ownership_requests.owner_id = ?", ownerId)

	if params.Search.City != "" {
		baseQuery = baseQuery.Where("properties.city ILIKE ?", "%"+params.Search.City+"%")
	}
	if params.Search.Street != "" {
		baseQuery = baseQuery.Where("properties.street ILIKE ?", "%"+params.Search.Street+"%")
	}
	if params.Search.Number != "" {
		baseQuery = baseQuery.Where("properties.number ILIKE ?", "%"+params.Search.Number+"%")
	}
	if params.Search.Floor != 0 {
		baseQuery = baseQuery.Where("households.floor = ?", params.Search.Floor)
	}
	if params.Search.Suite != "" {
		baseQuery = baseQuery.Where("households.suite = ?", "%"+params.Search.Suite+"%")
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		repository.Logger.Error("Error querying requests count", err)
		return nil, 0, err
	}

	sortOrder := params.SortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	query := baseQuery.Order(fmt.Sprintf("%s %s", params.SortBy, sortOrder))
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.
		Preload("Household").
		Preload("Household.Property").
		Preload("Owner").
		Find(&ownershipRequests).Error; err != nil {
		repository.Logger.Error("Error querying ownership requests", err)
		return nil, 0, err
	}

	return ownershipRequests, total, nil
}

func (repository *OwnershipRepository) FindById(id uint64) (*model.OwnershipRequest, error) {
	var ownershipRequest model.OwnershipRequest
	if err := repository.Database.First(&ownershipRequest, "id = ?", id).Error; err != nil {
		repository.Logger.Error("Error finding ownership request", err)
		return nil, err
	}
	return &ownershipRequest, nil
}

func (repository *OwnershipRepository) AcceptRequest(tx *gorm.DB, id uint64) error {
	const newRequestStatus model.RequestStatus = 1
	closedTime := time.Now()
	err := tx.Model(&model.OwnershipRequest{}).Where("id = ?", id).
		Update("status", newRequestStatus).
		Update("closed_time", closedTime).
		Error
	if err != nil {
		repository.Logger.Error("Error updating ownership request", err)
		return err
	}
	return nil
}
