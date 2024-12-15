package repository

import (
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
