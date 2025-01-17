package repository

import (
	"watt-flow/db"
	"watt-flow/model"
	"watt-flow/util"

	"gorm.io/gorm/clause"
)

type PricelistRepository struct {
	Database db.Database
	Logger   util.Logger
}

func NewPricelistRepository(db db.Database, logger util.Logger) *PricelistRepository {
	err := db.AutoMigrate(&model.TimeSlot{})
	if err != nil {
		logger.Error("Error migrating pricelist repo", err)
	}
	return &PricelistRepository{
		Database: db,
		Logger:   logger,
	}
}

func (repository *PricelistRepository) Create(pricelist *model.Pricelist) (model.Pricelist, error) {
	result := repository.Database.Preload(clause.Associations).Create(pricelist)
	if result.Error != nil {
		repository.Logger.Error("Error creating pricelist", result.Error)
		return *pricelist, result.Error
	}
	return *pricelist, nil
}

func (repository *PricelistRepository) FindById(id uint64) (*model.Pricelist, error) {
	var pricelist model.Pricelist
	if err := repository.Database.Preload(clause.Associations).Where("id = ?", id).First(&pricelist).Error; err != nil {
		repository.Logger.Error("Error finding pricelist by ID", err)
		return nil, err
	}
	return &pricelist, nil
}
