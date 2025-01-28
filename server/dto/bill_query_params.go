package dto

type BillSearchParams struct {
	// Clerk string `json:"clerk"`
}

type BillQueryParams struct {
	Page      int              `json:"page"`
	PageSize  int              `json:"pageSize"`
	SortBy    string           `json:"sortBy"`
	SortOrder string           `json:"sortOrder"`
	Search    BillSearchParams `json:"search"`
}
