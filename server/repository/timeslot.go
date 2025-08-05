package repository

import (
	"watt-flow/db"
	"watt-flow/model"
	"watt-flow/util"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TimeSlotRepository struct {
	Database db.Database
	Logger   util.Logger
}

func NewTimeSlotRepository(db db.Database, logger util.Logger) *TimeSlotRepository {
	err := db.AutoMigrate(&model.TimeSlot{})
	if err != nil {
		logger.Error("Error migrating timeslot repo", err)
	}
	return &TimeSlotRepository{
		Database: db,
		Logger:   logger,
	}
}

func (r *TimeSlotRepository) WithTrx(trxHandle *gorm.DB) *TimeSlotRepository {
	if trxHandle == nil {
		r.Logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	return &TimeSlotRepository{
		Database: db.Database{DB: trxHandle},
		Logger:   r.Logger,
	}
}

func (repository *TimeSlotRepository) DeleteSlotsForClerk(clerkId uint64) error {
	result := repository.Database.Where("clerk_id = ?", clerkId).Delete(&model.TimeSlot{})
	if result.Error != nil {
		repository.Logger.Error("Error deleting time slots for clerk", result.Error)
		return result.Error
	}
	//if result.RowsAffected == 0 {
	//	return fmt.Errorf("no slots found for clerk with ID %d", clerkId)
	//}
	return nil
}

func (repository *TimeSlotRepository) Create(timeslot *model.TimeSlot) (model.TimeSlot, error) {
	result := repository.Database.Preload(clause.Associations).Create(timeslot)
	if result.Error != nil {
		repository.Logger.Error("Error creating timeslot", result.Error)
		return *timeslot, result.Error
	}
	return *timeslot, nil
}

func (repository *TimeSlotRepository) FindByDate(date datatypes.Date) (*model.TimeSlot, error) {
	var timeslot model.TimeSlot
	if err := repository.Database.Preload(clause.Associations).Where("date = ?", date).First(&timeslot).Error; err != nil {
		repository.Logger.Error("Error finding timeslot by ID", err)
		return nil, err
	}
	return &timeslot, nil
}

func (repository *TimeSlotRepository) FindByDateAndClerkId(date datatypes.Date, clerkId uint64) (*model.TimeSlot, error) {
	var timeslot model.TimeSlot
	if err := repository.Database.Preload(clause.Associations).Where("date = ? AND clerk_id = ?", date, clerkId).First(&timeslot).Error; err != nil {
		repository.Logger.Error("Error finding timeslot by ID", err)
		return nil, err
	}
	return &timeslot, nil
}

func (repository *TimeSlotRepository) FindByDateAndClerkID(date datatypes.Date, clerkId uint64) (*model.TimeSlot, error) {
	var timeslot model.TimeSlot
	if err := repository.Database.Preload(clause.Associations).Where("date = ? AND clerk_id = ?", date, clerkId).First(&timeslot).Error; err != nil {
		repository.Logger.Error("Error finding timeslot by ID", err)
		return nil, err
	}
	return &timeslot, nil
}

func (repository *TimeSlotRepository) Update(timeslot *model.TimeSlot) (model.TimeSlot, error) {
	result := repository.Database.Save(timeslot)
	if result.Error != nil {
		repository.Logger.Error("Error updating timeslot", result.Error)
		return *timeslot, result.Error
	}
	return *timeslot, nil
}
