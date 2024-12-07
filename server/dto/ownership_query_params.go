package dto

type OwnershipSearchParams struct {
	City   string `json:"city"`
	Street string `json:"street"`
	Number string `json:"number"`
	Floor  int    `json:"floor"`
	Suite  string `json:"suite"`
}

type OwnershipQueryParams struct {
	Page      int                   `json:"page"`
	PageSize  int                   `json:"pageSize"`
	SortBy    string                `json:"sortBy"`
	SortOrder string                `json:"sortOrder"`
	Search    OwnershipSearchParams `json:"search"`
}
