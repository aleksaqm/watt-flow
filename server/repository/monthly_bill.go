package repository

import (
	"fmt"
	"time"
	"watt-flow/db"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MonthlyBillRepository struct {
	Database db.Database
	Logger   util.Logger
}

func NewMonthlyBillRepository(db db.Database, logger util.Logger) *MonthlyBillRepository {
	err := db.AutoMigrate(&model.MonthlyBill{})
	if err != nil {
		logger.Error("Error migrating monthly bill repo", err)
	}
	return &MonthlyBillRepository{
		Database: db,
		Logger:   logger,
	}
}

func (r *MonthlyBillRepository) WithTrx(trxHandle *gorm.DB) *MonthlyBillRepository {
	if trxHandle == nil {
		r.Logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	return &MonthlyBillRepository{
		Database: db.Database{DB: trxHandle},
		Logger:   r.Logger,
	}
}

func (repository *MonthlyBillRepository) Create(bill *model.MonthlyBill) (model.MonthlyBill, error) {
	result := repository.Database.Preload(clause.Associations).Create(bill)
	if result.Error != nil {
		repository.Logger.Error("Error creating monthly bill", result.Error)
		return *bill, result.Error
	}
	return *bill, nil
}

func (repository *MonthlyBillRepository) FindById(id uint64) (*model.MonthlyBill, error) {
	var bill model.MonthlyBill
	if err := repository.Database.Preload(clause.Associations).Where("id = ?", id).First(&bill).Error; err != nil {
		repository.Logger.Error("Error finding monthly bill by ID", err)
		return nil, err
	}
	return &bill, nil
}

func (repository *MonthlyBillRepository) FindBillsBetweenDates(date1 time.Time, date2 time.Time) ([]model.MonthlyBill, error) {
	var bills []model.MonthlyBill
	startKey := fmt.Sprintf("%d-%02d", date1.Year(), date1.Month())
	endKey := fmt.Sprintf("%d-%02d", date2.Year(), date2.Month())

	err := repository.Database.Where("billing_date BETWEEN ? AND ? AND status = 'Completed'", startKey, endKey).Find(&bills).Error
	if err != nil {
		return nil, err
	}
	return bills, nil
}

func (repository *MonthlyBillRepository) Query(params *dto.MonthlyBillQueryParams) ([]model.MonthlyBill, int64, error) {
	var bills []model.MonthlyBill
	var total int64

	baseQuery := repository.Database.Model(&model.MonthlyBill{})

	if err := baseQuery.Count(&total).Error; err != nil {
		repository.Logger.Error("Error querying bills count", err)
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
		Find(&bills).Error; err != nil {
		repository.Logger.Error("Error querying bills", err)
		return nil, 0, err
	}

	return bills, total, nil
}

func (repository *MonthlyBillRepository) Update(bill *model.MonthlyBill) (model.MonthlyBill, error) {
	result := repository.Database.Save(bill)
	if result.Error != nil {
		repository.Logger.Error("Error updating MonthlyBill", result.Error)
		return *bill, result.Error
	}
	return *bill, nil
}
