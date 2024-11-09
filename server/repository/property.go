package repository

import (
	"gorm.io/gorm/clause"
	"watt-flow/db"
	"watt-flow/model"
	"watt-flow/util"
)

type PropertyRepository struct {
	database db.Database
	logger   util.Logger
}

func NewPropertyRepository(db db.Database, logger util.Logger) *PropertyRepository {
	err := db.AutoMigrate(&model.Property{})
	if err != nil {
		logger.Error("Error migrating property", err)
	}
	return &PropertyRepository{
		database: db,
		logger:   logger,
	}
}

func (repository *PropertyRepository) Create(property *model.Property) (model.Property, error) {
	result := repository.database.Preload("Owner").Preload("Address").Create(property)
	if result.Error != nil {
		repository.logger.Error("Error creating property", result.Error)
		return *property, result.Error
	}
	return *property, nil
}

func (repository *PropertyRepository) FindById(id uint64) (*model.Property, error) {
	var property model.Property
	if err := repository.database.Preload(clause.Associations).First(&property, id).Error; err != nil {
		repository.logger.Error("Error finding property by ID", err)
		return nil, err
	}
	return &property, nil
}

func (repository *PropertyRepository) FindByStatus(status model.PropertyStatus) ([]model.Property, error) {
	var properties []model.Property
	result := repository.database.Where("status = ?", status).Find(&properties)
	if result.Error != nil {
		repository.logger.Error("Error finding properties by status", result.Error)
		return nil, result.Error
	}
	return properties, nil
}

func (repository *PropertyRepository) Update(property *model.Property) (model.Property, error) {
	result := repository.database.Save(property)
	if result.Error != nil {
		repository.logger.Error("Error updating property", result.Error)
		return *property, result.Error
	}
	return *property, nil
}

func (repository *PropertyRepository) Delete(id uint64) error {
	result := repository.database.Delete(&model.Property{}, id)
	if result.Error != nil {
		repository.logger.Error("Error deleting property", result.Error)
		return result.Error
	}
	return nil
}
