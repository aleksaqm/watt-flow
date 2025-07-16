package dto

import (
	"time"
	"watt-flow/model"
)

type BillTaskDto struct {
	BillingDate   string
	IssueDate     time.Time
	Pricelist     model.Pricelist
	OwnerID       uint64
	OwnerEmail    string
	OwnerUsername string
	PowerMeterID  string
}
