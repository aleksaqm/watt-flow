package dto

import "gorm.io/datatypes"

type UpdateTimeSlotDto struct {
	Date      datatypes.Date `json:"date"`
	ClerkId   uint64         `json:"clerkId"`
	Occupied  []int          `json:"occupied"`
	MeetingId uint64         `json:"meetingId"`
}
