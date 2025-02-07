package repository

import (
	"watt-flow/db"
	"watt-flow/model"
	"watt-flow/util"

	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
)

type TimeSlotRepository struct {
	Database db.Database
	Logger   util.Logger
}

func NewTimeSlotRepository(db db.Database, logger util.Logger) TimeSlotRepository {
	err := db.AutoMigrate(&model.TimeSlot{})
	if err != nil {
		logger.Error("Error migrating timeslot repo", err)
	}
	return TimeSlotRepository{
		Database: db,
		Logger:   logger,
	}
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
