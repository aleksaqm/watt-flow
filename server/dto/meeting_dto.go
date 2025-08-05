package dto

import "time"

type MeetingDTO struct {
	ID        uint64    `json:"id"`
	StartTime time.Time `json:"start_time"`
	Duration  int32     `json:"duration"`
	ClerkID   uint64    `json:"clerk_id"`
	UserID    uint64    `json:"user_id"`
	Username  string    `json:"username"`
}
