package model

import (
	"errors"
	"time"
)

type PropertyStatus int

const (
	PendingProperty PropertyStatus = iota
	DeniedProperty
	ApprovedProperty
)

func ParsePropertyStatus(status string) (PropertyStatus, error) {
	switch status {
	case "PendingProperty":
		return PendingProperty, nil
	case "DeniedProperty":
		return DeniedProperty, nil
	case "ApprovedProperty":
		return ApprovedProperty, nil
	default:
		return PendingProperty, errors.New("invalid status value")
	}
}

type Property struct {
	Id          uint64         `gorm:"primary_key;autoIncrement" json:"id"`
	Floors      int32          `json:"floors"`
	Images      []string       `gorm:"serializer:json" json:"images"`
	Documents   []string       `gorm:"serializer:json" json:"documents"`
	Status      PropertyStatus `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	ConfirmedAt time.Time      `json:"confirmed_at"`
	OwnerID     uint64         `gorm:"column:owner_id" json:"owner_id"`
	Address     Address        `gorm:"foreignKey:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;embedded" json:"address"`
	Owner       *User          `gorm:"foreignKey:OwnerID" json:"owner"`
	Household   []Household    `gorm:"foreignKey:PropertyID" json:"household"`
}
