package dto

import "gorm.io/datatypes"

type TimeSlotDto struct {
	Date    datatypes.Date `json:"date"`
	ClerkId uint64         `json:"clerkId"`
	Slots   [15]bool       `json:"slots"`
	Id      uint64         `json:"id"`
}
