package repository

import (
	"fmt"
	"watt-flow/db"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BillRepository struct {
	Database db.Database
	Logger   util.Logger
}

func NewBillRepository(db db.Database, logger util.Logger) BillRepository {
	err := db.AutoMigrate(&model.Bill{})
	if err != nil {
		logger.Error("Error migrating bill repo", err)
	}
	return BillRepository{
		Database: db,
		Logger:   logger,
	}
}

func (repository *BillRepository) Create(bill *model.Bill) (*model.Bill, error) {
	result := repository.Database.Preload(clause.Associations).Create(bill)
	if result.Error != nil {
		repository.Logger.Error("Error creating bill", result.Error)
		return bill, result.Error
	}
	return bill, nil
}

func (repository *BillRepository) FindById(id uint64) (*model.Bill, error) {
	var bill model.Bill
	if err := repository.Database.Preload(clause.Associations).Where("id = ?", id).First(&bill).Error; err != nil {
		repository.Logger.Error("Error finding bill by ID", err)
		return nil, err
	}
	return &bill, nil
}

func (r *BillRepository) UpdateStatusToPaid(billID uint64) error {
	result := r.Database.
		Model(&model.Bill{}).
		Where("id = ?", billID).
		Updates(map[string]interface{}{
			"status": "Paid",
		})

	if result.Error != nil {
		r.Logger.Error("Error updating bill status to paid", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r BillRepository) WithTrx(trxHandle *gorm.DB) BillRepository {
	if trxHandle == nil {
		r.Logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	r.Database.DB = trxHandle
	return r
}

func (repository *BillRepository) Update(bill *model.Bill) (model.Bill, error) {
	result := repository.Database.Save(bill)
	if result.Error != nil {
		repository.Logger.Error("Error updating Bill", result.Error)
		return *bill, result.Error
	}
	return *bill, nil
}

func (r *BillRepository) SearchBills(params *dto.BillQueryParams) ([]model.Bill, int64, error) {
	var bills []model.Bill
	var total int64

	db := r.Database.Model(&model.Bill{})

	searchParams := params.Search

	if searchParams.UserID != 0 {
		db = db.Where("owner_id = ?", searchParams.UserID)
	}
	if searchParams.Status != "" {
		db = db.Where("status = ?", searchParams.Status)
	}
	if searchParams.MinPrice > 0 {
		db = db.Where("price >= ?", searchParams.MinPrice)
	}
	if searchParams.MaxPrice > 0 {
		db = db.Where("price <= ?", searchParams.MaxPrice)
	}
	if searchParams.BillingDate != "" {
		db = db.Where("billing_date = ?", searchParams.BillingDate)
	}

	if err := db.Count(&total).Error; err != nil {
		r.Logger.Error("Error counting bills", err)
		return nil, 0, err
	}

	orderByClause := fmt.Sprintf("%s %s", params.SortBy, params.SortOrder)
	db = db.Order(orderByClause)

	offset := (params.Page - 1) * params.PageSize
	db = db.Offset(offset).Limit(params.PageSize)

	if err := db.Preload("Owner").Preload("Pricelist").Find(&bills).Error; err != nil {
		r.Logger.Error("Error finding bills with query", err)
		return nil, 0, err
	}

	return bills, total, nil
}

// func (repository *BillRepository) FindForUser(userID uint64, params *dto.BillSearchParams) ([]model.Bill, int64, error) {
// 	var bills []model.Bill
// 	var total int64
//
// 	baseQuery := repository.Database.Model(&model.Bill{}).
// 		Joins("JOIN users ON users.id = meetings.clerk_id").
// 		Where("meetings.user_id = ?", userID)
//
// 	if params.Search.Clerk != "" {
// 		baseQuery = baseQuery.Where("users.username ILIKE ?", "%"+params.Search.Clerk+"%")
// 	}
//
// 	if err := baseQuery.Count(&total).Error; err != nil {
// 		repository.Logger.Error("Error querying meetings count", err)
// 		return nil, 0, err
// 	}
//
// 	sortOrder := params.SortOrder
// 	if sortOrder != "asc" && sortOrder != "desc" {
// 		sortOrder = "asc"
// 	}
//
// 	query := baseQuery.Order(fmt.Sprintf("%s %s", params.SortBy, sortOrder))
// 	offset := (params.Page - 1) * params.PageSize
// 	query = query.Offset(offset).Limit(params.PageSize)
//
// 	if err := query.
// 		Preload("Clerk").
// 		Find(&meetings).Error; err != nil {
// 		repository.Logger.Error("Error querying meetings", err)
// 		return nil, 0, err
// 	}
//
// 	return meetings, total, nil
// }
