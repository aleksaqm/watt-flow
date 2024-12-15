package repository

import (
	"gorm.io/gorm/clause"
	"watt-flow/db"
	"watt-flow/model"
	"watt-flow/util"
)

type MeetingRepository struct {
	Database db.Database
	Logger   util.Logger
}

func NewMeetingRepository(db db.Database, logger util.Logger) *MeetingRepository {
	err := db.AutoMigrate(&model.Meeting{})
	if err != nil {
		logger.Error("Error migrating meeting repo", err)
	}
	return &MeetingRepository{
		Database: db,
		Logger:   logger,
	}
}

func (repository *MeetingRepository) Create(meeting *model.Meeting) (model.Meeting, error) {
	result := repository.Database.Preload(clause.Associations).Create(meeting)
	if result.Error != nil {
		repository.Logger.Error("Error creating meeting", result.Error)
		return *meeting, result.Error
	}
	return *meeting, nil
}
func (repository *MeetingRepository) FindById(id uint64) (*model.Meeting, error) {
	var meeting model.Meeting
	if err := repository.Database.Preload(clause.Associations).Where("id = ?", id).First(&meeting).Error; err != nil {
		repository.Logger.Error("Error finding meeting by ID", err)
		return nil, err
	}
	return &meeting, nil
}
func (repository *MeetingRepository) FindBySlotId(id uint64) (*model.Meeting, error) {
	var meeting model.Meeting
	if err := repository.Database.Preload(clause.Associations).Where("time_slot_id = ?", id).First(&meeting).Error; err != nil {
		repository.Logger.Error("Error finding meeting by ID", err)
		return nil, err
	}
	return &meeting, nil
}
