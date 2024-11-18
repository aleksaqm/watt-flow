package dto

type PropertySearchParams struct {
	City   string `json:"city"`
	Street string `json:"street"`
	Number string `json:"number"`
	Floors int    `json:"floors"`
}

type PropertyQueryParams struct {
	Page      int                  `json:"page"`
	PageSize  int                  `json:"pageSize"`
	SortBy    string               `json:"sortBy"`
	SortOrder string               `json:"sortOrder"`
	Search    PropertySearchParams `json:"search"`
}
