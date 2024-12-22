package dto

type NewMeetingDTO struct {
	Meeting  MeetingDTO        `json:"meeting"`
	TimeSlot UpdateTimeSlotDto `json:"timeslot"`
}
