package repository

import (
	"fmt"
	"gorm.io/gorm/clause"
	"watt-flow/db"
	"watt-flow/dto"
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

func (repository *MeetingRepository) FindForUser(userID uint64, params *dto.MeetingQueryParams) ([]model.Meeting, int64, error) {
	var meetings []model.Meeting
	var total int64

	baseQuery := repository.Database.Model(&model.Meeting{}).
		Joins("JOIN users ON users.id = meetings.clerk_id").
		Where("meetings.user_id = ?", userID)

	if params.Search.Clerk != "" {
		baseQuery = baseQuery.Where("users.username ILIKE ?", "%"+params.Search.Clerk+"%")
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		repository.Logger.Error("Error querying meetings count", err)
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
		Preload("Clerk").
		Find(&meetings).Error; err != nil {
		repository.Logger.Error("Error querying meetings", err)
		return nil, 0, err
	}

	return meetings, total, nil

}
