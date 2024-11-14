package repository

import (
	"gorm.io/gorm/clause"
	"watt-flow/db"
	"watt-flow/model"
	"watt-flow/util"
)

type HouseholdRepository struct {
	Database db.Database
	Logger   util.Logger
}

func NewHouseholdRepository(db db.Database, logger util.Logger) *HouseholdRepository {
	err := db.AutoMigrate(&model.Household{})
	if err != nil {
		logger.Error("Error migrating household", err)
	}
	return &HouseholdRepository{
		Database: db,
		Logger:   logger,
	}
}

func (repository *HouseholdRepository) Create(household *model.Household) (model.Household, error) {
	result := repository.Database.Preload(clause.Associations).Create(household)
	if result.Error != nil {
		repository.Logger.Error("Error creating household", result.Error)
		return *household, result.Error
	}
	return *household, nil
}

func (repository *HouseholdRepository) FindById(id uint64) (*model.Household, error) {
	var household model.Household
	if err := repository.Database.Preload(clause.Associations).First(&household, id).Error; err != nil {
		repository.Logger.Error("Error finding household by ID", err)
		return nil, err
	}
	return &household, nil
}

func (repository *HouseholdRepository) FindByStatus(status model.HouseholdStatus) ([]model.Household, error) {
	var households []model.Household
	result := repository.Database.Where("status = ?", status).Find(&households)
	if result.Error != nil {
		repository.Logger.Error("Error finding households by status", result.Error)
		return nil, result.Error
	}
	return households, nil
}

func (repository *HouseholdRepository) Update(household *model.Household) (model.Household, error) {
	result := repository.Database.Save(household)
	if result.Error != nil {
		repository.Logger.Error("Error updating household", result.Error)
		return *household, result.Error
	}
	return *household, nil
}

func (repository *HouseholdRepository) Delete(id uint64) error {
	result := repository.Database.Delete(&model.Household{}, id)
	if result.Error != nil {
		repository.Logger.Error("Error deleting household", result.Error)
		return result.Error
	}
	return nil
}
