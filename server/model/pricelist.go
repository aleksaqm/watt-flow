package model

import (
	"gorm.io/datatypes"
)

type Pricelist struct {
	ID           uint64         `gorm:"primary_key" json:"id"`
	ValidFrom    datatypes.Date `json:"valid_from"`
	BlueZone     float64        `json:"blue_zone"`
	RedZone      float64        `json:"red_zone"`
	GreenZone    float64        `json:"green_zone"`
	BillingPower float64        `json:"billing_power"`
	Tax          float64        `json:"tax"`
	IsActive     bool           `json:"is_active"`
}
