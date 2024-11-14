package repository

import (
	"fmt"
	"gorm.io/gorm/clause"
	"watt-flow/db"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/util"
)

type PropertyRepository struct {
	Database db.Database
	Logger   util.Logger
}

func NewPropertyRepository(db db.Database, logger util.Logger) *PropertyRepository {
	err := db.AutoMigrate(&model.Property{})
	if err != nil {
		logger.Error("Error migrating property", err)
	}
	return &PropertyRepository{
		Database: db,
		Logger:   logger,
	}
}

func (repository *PropertyRepository) Create(property *model.Property) (model.Property, error) {
	result := repository.Database.Create(property)
	if result.Error != nil {
		repository.Logger.Error("Error creating property", result.Error)
		return *property, result.Error
	}
	return *property, nil
}

func (repository *PropertyRepository) FindById(id uint64) (*model.Property, error) {
	var property model.Property
	if err := repository.Database.Preload(clause.Associations).First(&property, id).Error; err != nil {
		repository.Logger.Error("Error finding property by ID", err)
		return nil, err
	}
	return &property, nil
}

func (repository *PropertyRepository) FindByStatus(status model.PropertyStatus) ([]model.Property, error) {
	var properties []model.Property
	result := repository.Database.Where("status = ?", status).Find(&properties)
	if result.Error != nil {
		repository.Logger.Error("Error finding properties by status", result.Error)
		return nil, result.Error
	}
	return properties, nil
}

func (repository *PropertyRepository) Update(property *model.Property) (model.Property, error) {
	result := repository.Database.Save(property)
	if result.Error != nil {
		repository.Logger.Error("Error updating property", result.Error)
		return *property, result.Error
	}
	return *property, nil
}

func (repository *PropertyRepository) Delete(id uint64) error {
	result := repository.Database.Delete(&model.Property{}, id)
	if result.Error != nil {
		repository.Logger.Error("Error deleting property", result.Error)
		return result.Error
	}
	return nil
}

func (repository *PropertyRepository) TableQuery(params *dto.PropertyQueryParams) ([]model.Property, int64, error) {
	var properties []model.Property
	var total int64

	baseQuery := repository.Database.Model(&model.Property{}).
		Preload("Owner").
		Preload("Household")

	if params.Search != "" {
		baseQuery = baseQuery.Where("city ILIKE ? OR street ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		repository.Logger.Error("Error querying property count", err)
		return nil, 0, err
	}

	sortOrder := params.SortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	query := baseQuery.Order(fmt.Sprintf("%s %s", params.SortBy, sortOrder))

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&properties).Error; err != nil {
		repository.Logger.Error("Error querying property", err)
		return nil, 0, err
	}

	return properties, total, nil
}
