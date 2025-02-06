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
}

type Pricelist struct {
	ID           uint64
	ValidFrom    time.Time
	BlueZone     float64
	RedZone      float64
	GreenZone    float64
	BillingPower float64
	Tax          float64
}
type Bill struct {
	BillingDate string
	IssueDate   time.Time
	PricelistID uint64
	OwnerID     uint64
	SpentPower  float64
	Price       float64
	Status      string
}
