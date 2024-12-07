package dto

type UserSearchParams struct {
	Id        uint64 `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
}

type UserQueryParams struct {
	Page      int              `json:"page"`
	PageSize  int              `json:"pageSize"`
	SortBy    string           `json:"sortBy"`
	SortOrder string           `json:"sortOrder"`
	Search    UserSearchParams `json:"search"`
}
