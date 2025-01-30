package model

import "time"

type MonthlyBill struct {
	ID          uint64    `gorm:"primary_key" json:"id"`
	IssueDate   time.Time `json:"issue_date"`
	BillingDate string    `gorm:"unique" json:"billing_date"`
	Status      string    `json:"status"`
}
