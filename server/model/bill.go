package model

import (
	"gorm.io/datatypes"
)

type Bill struct {
	ID          uint64         `gorm:"primary_key" json:"id"`
	Date        datatypes.Date `json:"date"`
	Pricelist   Pricelist      `gorm:"foreignKey:PricelistID" json:"pricelist"`
	PricelistID uint64         `json:"pricelist_id"`
	SpentPower  float64        `json:"spent_power"`
	Price       float64        `json:"price"`
	Owner       User           `gorm:"foreignKey:OwnerID" json:"owner"`
	OwnerID     uint64         `json:"owner_id"`
}
