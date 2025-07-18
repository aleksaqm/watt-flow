package dto

type GrantAccessRequestDto struct {
	UserID uint64 `json:"userId" binding:"required"`
}
