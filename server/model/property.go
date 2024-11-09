package model

import "time"

type PropertyStatus int

const (
	PendingProperty PropertyStatus = iota
	DeniedProperty
	ApprovedProperty
)

type Property struct {
	Id          uint64         `gorm:"primary_key" json:"id"`
	Floors      int32          `json:"floors"`
	Images      []string       `gorm:"serializer:json" json:"images"`
	Documents   []string       `gorm:"serializer:json" json:"documents"`
	Status      PropertyStatus `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	ConfirmedAt time.Time      `json:"confirmed_at"`
	AddressID   uint64         `gorm:"column:address_id" json:"address_id"`
	OwnerID     uint64         `gorm:"column:owner_id" json:"owner_id"`
	Address     *Address       `gorm:"foreignKey:AddressID" json:"address"`
	Owner       *User          `gorm:"foreignKey:OwnerID" json:"owner"`
	Household   []Household    `gorm:"foreignKey:PropertyID" json:"household"`
}
