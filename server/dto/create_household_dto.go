package dto

type CreateHouseholdDto struct {
	Floor           uint32  `json:"floor"`
	Suite           string  `json:"suite"`
	SqFootage       float32 `json:"sq_footage"`
	PropertyId      uint64  `json:"property_id"`
	CadastralNumber string  `json:"cadastral_number"`
}
