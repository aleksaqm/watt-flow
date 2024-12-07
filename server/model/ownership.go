package model

import "time"

type RequestStatus int

const (
	Pending RequestStatus = iota
	Approved
	Declined
)

type OwnershipRequest struct {
	Id           uint64        `gorm:"primaryKey;autoIncrement" json:"id"`
	Images       []string      `gorm:"serializer:json" json:"images"`
	Documents    []string      `gorm:"serializer:json" json:"documents"`
	Status       RequestStatus `json:"status"`
	DenialReason string        `json:"denial_reason"`
	OwnerID      uint64        `gorm:"column:owner_id;null" json:"owner_id"`
	Owner        *User         `gorm:"foreignKey:OwnerID" json:"owner"`
	HouseholdID  uint64        `gorm:"column:household_id;null" json:"household_id"`
	Household    *Household    `gorm:"foreignKey:HouseholdID" json:"household"`
	CreatedAt    time.Time     `json:"created_at"`
	ClosedAt     time.Time     `json:"closed_at"`
}
