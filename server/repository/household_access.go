package repository

import (
	"watt-flow/db"
	"watt-flow/model"
	"watt-flow/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HouseholdAccessRepository struct {
	Database db.Database
	Logger   util.Logger
}

func NewHouseholdAccessRepository(db db.Database, logger util.Logger) *HouseholdAccessRepository {
	err := db.AutoMigrate(&model.HouseholdAccess{})
	if err != nil {
		logger.Error("Error migrating household_access", err)
	}
	return &HouseholdAccessRepository{
		Database: db,
		Logger:   logger,
	}
}

func (r *HouseholdAccessRepository) WithTrx(trxHandle *gorm.DB) *HouseholdAccessRepository {
	if trxHandle == nil {
		r.Logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	return &HouseholdAccessRepository{
		Database: db.Database{DB: trxHandle},
		Logger:   r.Logger,
	}
}

func (r HouseholdAccessRepository) Create(access *model.HouseholdAccess) error {
	result := r.Database.Create(access)
	if result.Error != nil {
		r.Logger.Error("Error creating household access", result.Error)
		return result.Error
	}
	return nil
}

func (r HouseholdAccessRepository) FindByHouseholdIDAndUserID(householdID, userID uint64) (*model.HouseholdAccess, error) {
	var access model.HouseholdAccess
	result := r.Database.Where("household_id = ? AND user_id = ?", householdID, userID).First(&access)
	if result.Error != nil {
		return nil, result.Error
	}
	r.Logger.Info("Existing household access: ", &access)
	return &access, nil
}

func (r HouseholdAccessRepository) FindByHouseholdID(householdID uint64) ([]model.HouseholdAccess, error) {
	var access []model.HouseholdAccess
	result := r.Database.Preload(clause.Associations).Where("household_id = ?", householdID).Find(&access)
	if result.Error != nil {
		return nil, result.Error
	}
	r.Logger.Info("Existing user with access: ", access[0].User.Id, access[0].User.Username)
	r.Logger.Info("Existing household access: ", &access)
	return access, nil
}

func (r HouseholdAccessRepository) Delete(householdID, userID uint64) error {
	result := r.Database.Where("household_id = ? AND user_id = ?", householdID, userID).Delete(&model.HouseholdAccess{})
	if result.Error != nil {
		r.Logger.Error("Error deleting household access", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
