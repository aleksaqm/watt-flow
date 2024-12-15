package model

import "time"

type Meeting struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	StartTime time.Time `json:"start_time"`
	Duration  int32     `json:"duration"`
	ClerkID   uint64    `json:"clerk_id"`
	UserID    uint64    `json:"user_id"`
	Clerk     User      `gorm:"foreignKey:ClerkID" json:"clerk"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}
