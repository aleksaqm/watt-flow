package dto

type OwnershipRequestDto struct {
	Id          uint64   `json:"id"`
	Images      []string `json:"images"`
	Documents   []string `json:"documents"`
	OwnerID     uint64   `json:"owner_id"`
	HouseholdID uint64   `json:"household_id"`
}
