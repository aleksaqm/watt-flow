package dto

type CreatePropertyDto struct {
	Floors    int32    `json:"floors"`
	Images    []string `json:"images"`
	Documents []string `json:"documents"`
	AddressID uint64   `json:"address_id"`
}
