package dto

type OwnershipResponseDto struct {
	Id          uint64   `json:"id"`
	Images      []string `json:"images"`
	Documents   []string `json:"documents"`
	OwnerID     uint64   `json:"owner_id"`
	HouseholdID uint64   `json:"household_id"`
	CreatedAt   string   `json:"created_at"`
	ClosedAt    string   `json:"closed_at"`
	City        string   `json:"city"`
	Street      string   `json:"street"`
	Number      string   `json:"number"`
	Floor       uint32   `json:"floor"`
	Suite       string   `json:"suite"`
	Username    string   `json:"username"`
	Status      string   `json:"status"`
}
