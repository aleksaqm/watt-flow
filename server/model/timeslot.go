package model

import (
	"gorm.io/datatypes"
)

type TimeSlot struct {
	Date    datatypes.Date `gorm:"unique;type:date"`
	Clerk   User           `gorm:"foreignKey:ClerkID" json:"clerk"`
	ClerkID uint64         `json:"clerkId"`
	Slots   datatypes.JSON `gorm:"type:jsonb;not null"`
	Id      uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
}
