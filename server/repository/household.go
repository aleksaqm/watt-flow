package repository

import (
	"errors"
	"fmt"
	"watt-flow/db"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/util"

	"gorm.io/gorm"

	"gorm.io/gorm/clause"
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

func (r *HouseholdRepository) WithTrx(trxHandle *gorm.DB) *HouseholdRepository {
	if trxHandle == nil {
		r.Logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	return &HouseholdRepository{
		Database: db.Database{DB: trxHandle},
		Logger:   r.Logger,
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

func (repository *HouseholdRepository) FindMyHouseholdById(id uint64, userId uint64) (*model.Household, error) {
	var household model.Household

	db := repository.Database.
		Preload(clause.Associations).
		Where("id = ?", id).
		Where(`
            owner_id = ? OR
            EXISTS (
                SELECT 1 FROM household_accesses ha
                WHERE ha.household_id = households.id AND ha.user_id = ?
            )
        `, userId, userId)

	err := db.First(&household).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			repository.Logger.Warn("Household not found or user does not have access",
				map[string]interface{}{"householdID": id, "userID": userId})
			return nil, gorm.ErrRecordNotFound
		}
		repository.Logger.Error("Error finding household by ID with access check", err)
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

func (repository *HouseholdRepository) GetOwnedHouseholds() ([]model.Household, error) {
    var allHouseholds []model.Household
    batchSize := 10000
    offset := 0

    for {
        var batch []model.Household
        result := repository.Database.
            Where("status = 1").
            Preload("Owner").
            Limit(batchSize).
            Offset(offset).
            Find(&batch)

        if result.Error != nil {
            repository.Logger.Error("Error finding households by status", result.Error)
            return nil, result.Error
        }

        if len(batch) == 0 {
            break
        }

        allHouseholds = append(allHouseholds, batch...)
        offset += batchSize
    }

    return allHouseholds, nil
}

func (repository *HouseholdRepository) Query(params *dto.HouseholdQueryParams) ([]model.Household, int64, error) {
	var households []model.Household
	var total int64

	baseQuery := repository.Database.Model(&model.Household{})

	baseQuery = baseQuery.Joins("JOIN properties ON properties.id = households.property_id")

	repository.Logger.Info("AAAAAAA")
	repository.Logger.Info(params.Search.WithoutOwner)

	if params.Search.WithoutOwner {
		baseQuery = baseQuery.Where("households.status = ?", 2)
	}

	if params.Search.OwnerId != "" {
		baseQuery = baseQuery.Where(`
            households.owner_id = ? OR
            EXISTS (
                SELECT 1 FROM household_accesses ha
                WHERE ha.household_id = households.id AND ha.user_id = ?
            )
        `, params.Search.OwnerId, params.Search.OwnerId)
	}

	if params.Search.City != "" {
		baseQuery = baseQuery.Where("city ILIKE ?", "%"+params.Search.City+"%")
	}
	if params.Search.Street != "" {
		baseQuery = baseQuery.Where("street ILIKE ?", "%"+params.Search.Street+"%")
	}
	if params.Search.Number != "" {
		baseQuery = baseQuery.Where("number ILIKE ?", "%"+params.Search.Number+"%")
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

	if err := query.
		Preload("Property").
		Preload("Owner").
		Preload("DeviceStatus").
		Find(&households).Error; err != nil {
		repository.Logger.Error("Error querying households", err)
		return nil, 0, err
	}

	return households, total, nil
}

func (repository *HouseholdRepository) FindByCadastralNumber(cadastralNumber string) (*model.Household, error) {
	var household model.Household
	result := repository.Database.Where("cadastral_number = ?", cadastralNumber).Preload("Owner").Preload("Property").Preload("DeviceStatus").First(&household)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			repository.Logger.Warn("No household found for cadastralNumber", "cadastralNumber", cadastralNumber)
			return nil, nil
		}
		repository.Logger.Error("Error finding household by cadastralNumber", result.Error)
		return nil, result.Error
	}
	return &household, nil
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

func (repository *HouseholdRepository) AcceptHouseholds(tx *gorm.DB, propertyID uint64) error {
	const newHouseholdStatus model.HouseholdStatus = 2

	err := tx.Model(&model.Household{}).
		Where("property_id = ?", propertyID).
		Update("status", newHouseholdStatus).
		Error
	if err != nil {
		repository.Logger.Error("Error updating household status", err)
		return err
	}

	return nil
}

func (repository *HouseholdRepository) AddOwnerToHousehold(tx *gorm.DB, householdId uint64, ownerId uint64) error {
	err := tx.Model(&model.Household{}).
		Where("id = ?", householdId).
		Update("owner_id", ownerId).
		Update("status", 1).
		Error
	if err != nil {
		repository.Logger.Error("Error updating household owner", err)
		return err
	}
	return nil
}
