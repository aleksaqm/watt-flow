package model

import (
	"gorm.io/datatypes"
	"gorm.io/plugin/optimisticlock"
)

type TimeSlot struct {
	Date    datatypes.Date `gorm:"type:date"` //not unique
	Clerk   User           `gorm:"foreignKey:ClerkID" json:"clerk"`
	ClerkID uint64         `json:"clerkId"`
	Slots   datatypes.JSON `gorm:"type:jsonb;not null"`
	Id      uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Version optimisticlock.Version
}
