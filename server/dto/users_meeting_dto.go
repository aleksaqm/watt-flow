package dto

import "time"

type UsersMeetingDTO struct {
	ID        uint64    `json:"id"`
	StartTime time.Time `json:"start_time"`
	Duration  int32     `json:"duration"`
	Clerk     string    `json:"clerk"`
}
