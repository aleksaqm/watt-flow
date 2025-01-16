package dto

type MeetingSearchParams struct {
	Clerk string `json:"clerk"`
}

type MeetingQueryParams struct {
	Page      int                 `json:"page"`
	PageSize  int                 `json:"pageSize"`
	SortBy    string              `json:"sortBy"`
	SortOrder string              `json:"sortOrder"`
	Search    MeetingSearchParams `json:"search"`
}
