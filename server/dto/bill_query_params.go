package dto

import (
	"gorm.io/datatypes"
	"time"
)

type BillSearchParams struct {
	UserID           uint64  `json:"userId"`
	MinPrice         float64 `json:"minPrice"`
	MaxPrice         float64 `json:"maxPrice"`
	Status           string  `json:"status"`
	BillingDate      string  `json:"billingDate"`
	HouseholdID      uint64  `json:"householdId"`
	PaymentReference string  `json:"paymentReference"`
}

type BillQueryParams struct {
	Page      int              `json:"page"`
	PageSize  int              `json:"pageSize"`
	SortBy    string           `json:"sortBy"`
	SortOrder string           `json:"sortOrder"`
	Search    BillSearchParams `json:"search"`
}

type OwnerDto struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type PricelistDto struct {
	ID           uint64         `json:"id"`
	ValidFrom    datatypes.Date `json:"valid_from"`
	BlueZone     float64        `json:"blue_zone"`
	RedZone      float64        `json:"red_zone"`
	GreenZone    float64        `json:"green_zone"`
	BillingPower float64        `json:"billing_power"`
	Tax          float64        `json:"tax"`
}

type BillResponseDto struct {
	ID               uint64             `json:"id"`
	IssueDate        time.Time          `json:"issue_date"`
	BillingDate      string             `json:"billing_date"`
	SpentPower       float64            `json:"spent_power"`
	Price            float64            `json:"price"`
	Status           string             `json:"status"`
	Pricelist        PricelistDto       `json:"pricelist"`
	Owner            OwnerDto           `json:"owner"`
	Household        HouseholdResultDto `json:"household"`
	PaymentReference string             `json:"payment_reference"`
}
