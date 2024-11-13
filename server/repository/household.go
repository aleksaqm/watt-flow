package repository

import (
	"fmt"
	"watt-flow/db"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/util"

	"gorm.io/gorm/clause"
)

type HouseholdRepository struct {
	database db.Database
	logger   util.Logger
}

func NewHouseholdRepository(db db.Database, logger util.Logger) *HouseholdRepository {
	err := db.AutoMigrate(&model.Household{})
	if err != nil {
		logger.Error("Error migrating household", err)
	}
	return &HouseholdRepository{
		database: db,
		logger:   logger,
	}
}

func (repository *HouseholdRepository) Create(household *model.Household) (model.Household, error) {
	result := repository.database.Preload(clause.Associations).Create(household)
	if result.Error != nil {
		repository.logger.Error("Error creating household", result.Error)
		return *household, result.Error
	}
	return *household, nil
}

func (repository *HouseholdRepository) FindById(id uint64) (*model.Household, error) {
	var household model.Household
	if err := repository.database.Preload(clause.Associations).First(&household, id).Error; err != nil {
		repository.logger.Error("Error finding household by ID", err)
		return nil, err
	}
	return &household, nil
}

func (repository *HouseholdRepository) FindByStatus(status model.HouseholdStatus) ([]model.Household, error) {
	var households []model.Household
	result := repository.database.Where("status = ?", status).Find(&households)
	if result.Error != nil {
		repository.logger.Error("Error finding households by status", result.Error)
		return nil, result.Error
	}
	return households, nil
}

func (repository *HouseholdRepository) Search(searchDto dto.SearchHouseholdDto) ([]model.Household, error) {
	var households []model.Household

	query := repository.database.Model(&model.Property{}).Joins("Property").Joins("Property.Address")
	if searchDto.City != "" {
		query = query.Where("Property.Address.city = ?", searchDto.City)
	}
	if searchDto.Street != "" {
		query = query.Where("Property.Address.street LIKE ?", "%"+searchDto.Street+"%")
	}
	if searchDto.Number != "" {
		query = query.Where("Property.Address.number = ?", searchDto.Number)
	}
	if searchDto.Floor != 0 {
		query = query.Where("floor = ?", searchDto.Floor)
	}

	if searchDto.Suite != "" {
		query = query.Where("suite = ?", searchDto.Suite)
	}

	if err := query.Find(&households).Error; err != nil {
		fmt.Println("Error finding households:", err)
		return households, err
	}
	return households, nil
}

func (repository *HouseholdRepository) FindByCadastralNumber(cadastralNumber string) (*model.Household, error) {
	var household model.Household
	result := repository.database.Where("cadastral_number = ?", cadastralNumber).Find(&household)
	if result.Error != nil {
		repository.logger.Error("Error finding households by cadastralNumber", result.Error)
		return nil, result.Error
	}
	return &household, nil
}

func (repository *HouseholdRepository) Update(household *model.Household) (model.Household, error) {
	result := repository.database.Save(household)
	if result.Error != nil {
		repository.logger.Error("Error updating household", result.Error)
		return *household, result.Error
	}
	return *household, nil
}

func (repository *HouseholdRepository) Delete(id uint64) error {
	result := repository.database.Delete(&model.Household{}, id)
	if result.Error != nil {
		repository.logger.Error("Error deleting household", result.Error)
		return result.Error
	}
	return nil
}
