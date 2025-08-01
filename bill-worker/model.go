package main

import "time"

type BillTaskDto struct {
	BillingDate   string
	IssueDate     time.Time
	Pricelist     Pricelist
	OwnerID       uint64
	OwnerEmail    string
	OwnerUsername string
	PowerMeterID  string
	Last          bool
	MonthlyBillID uint64
	HouseHoldID   uint64
	HouseholdCN   string // Household Cadastral number
}

type Pricelist struct {
	ID           uint64    `json:"id"`
	ValidFrom    time.Time `json:"valid_from"`
	BlueZone     float64   `json:"blue_zone"`
	RedZone      float64   `json:"red_zone"`
	GreenZone    float64   `json:"green_zone"`
	BillingPower float64   `json:"billing_power"`
	Tax          float64   `json:"tax"`
}
type Bill struct {
	BillingDate      string
	IssueDate        time.Time
	PricelistID      uint64
	OwnerID          uint64
	SpentPower       float64
	Price            float64
	Status           string
	HouseholdID      uint64
	PaymentReference string
}
