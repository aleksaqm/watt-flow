package model

import "time"

type Bill struct {
	ID          uint64    `gorm:"primary_key" json:"id"`
	IssueDate   time.Time `json:"issue_date"`
	BillingDate string    `json:"billing_date"`
	Pricelist   Pricelist `gorm:"foreignKey:PricelistID" json:"pricelist"`
	PricelistID uint64    `json:"pricelist_id"`
	SpentPower  float64   `json:"spent_power"`
	Price       float64   `json:"price"`
	Owner       User      `gorm:"foreignKey:OwnerID" json:"owner"`
	OwnerID     uint64    `json:"owner_id"`
	Status      string    `json:"status"`
	HouseholdID uint64    `json:"household_id"`
	Household   Household `gorm:"foreignKey:HouseholdID"`
}
