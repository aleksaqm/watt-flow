package dto

type HouseholdSearchParams struct {
	Street       string `json:"street"`
	Number       string `json:"number"`
	City         string `json:"city"`
	Id           string `json:"id"`
	WithoutOwner bool   `json:"without_owner"`
}

type HouseholdQueryParams struct {
	Page      int                   `json:"page"`
	PageSize  int                   `json:"pageSize"`
	SortBy    string                `json:"sortBy"`
	SortOrder string                `json:"sortOrder"`
	Search    HouseholdSearchParams `json:"search"`
}
