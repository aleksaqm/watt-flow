package repository

import (
	"fmt"
	"time"
	"watt-flow/db"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/util"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PricelistRepository struct {
	Database db.Database
	Logger   util.Logger
}

func NewPricelistRepository(db db.Database, logger util.Logger) *PricelistRepository {
	err := db.AutoMigrate(&model.Pricelist{})
	if err != nil {
		logger.Error("Error migrating pricelist repo", err)
	}
	return &PricelistRepository{
		Database: db,
		Logger:   logger,
	}
}

func (r *PricelistRepository) WithTrx(trxHandle *gorm.DB) *PricelistRepository {
	if trxHandle == nil {
		r.Logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	return &PricelistRepository{
		Database: db.Database{DB: trxHandle},
		Logger:   r.Logger,
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

func (repository *PricelistRepository) GetActivePricelist() (*model.Pricelist, error) {
	var currentPricelist model.Pricelist
	err := repository.Database.Where("valid_from <= ?", time.Now()).
		Order("valid_from DESC").
		First(&currentPricelist).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching the current pricelist: %v", err)
	}
	return &currentPricelist, nil
}

func (repository *PricelistRepository) FindById(id uint64) (*model.Pricelist, error) {
	var pricelist model.Pricelist
	if err := repository.Database.Preload(clause.Associations).Where("id = ?", id).First(&pricelist).Error; err != nil {
		repository.Logger.Error("Error finding pricelist by ID", err)
		return nil, err
	}
	return &pricelist, nil
}

func (repository *PricelistRepository) FindByDate(date datatypes.Date) (*model.Pricelist, error) {
	var pricelist model.Pricelist
	if err := repository.Database.Preload(clause.Associations).Where("valid_from = ?", date).First(&pricelist).Error; err != nil {
		repository.Logger.Error("Error finding pricelist by ID", err)
		return nil, err
	}
	return &pricelist, nil
}

func (repository *PricelistRepository) Query(params *dto.PricelistQueryParams) ([]model.Pricelist, int64, error) {
	var pricelists []model.Pricelist
	var total int64

	baseQuery := repository.Database.Model(&model.Pricelist{})

	if err := baseQuery.Count(&total).Error; err != nil {
		repository.Logger.Error("Error querying pricelist count", err)
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
		Find(&pricelists).Error; err != nil {
		repository.Logger.Error("Error querying pricelists", err)
		return nil, 0, err
	}

	return pricelists, total, nil
}

func (repository *PricelistRepository) Delete(id uint64) error {
	result := repository.Database.Delete(&model.Pricelist{}, id)
	if result.Error != nil {
		repository.Logger.Error("Error deleting pricelist", result.Error)
		return result.Error
	}
	return nil
}
