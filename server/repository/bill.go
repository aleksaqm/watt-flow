package repository

import (
	"watt-flow/db"
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

func (r BillRepository) WithTrx(trxHandle *gorm.DB) BillRepository {
	if trxHandle == nil {
		r.Logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	r.Database.DB = trxHandle
	return r
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
