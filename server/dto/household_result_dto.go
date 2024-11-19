package dto

type HouseholdResultDto struct {
	Id              uint64  `json:"id"`
	Floor           uint32  `json:"floor"`
	Suite           string  `json:"suite"`
	Status          string  `json:"status"`
	SqFootage       float32 `json:"sq_footage"`
	OwnerID         uint64  `json:"owner_id"`
	OwnerName       string  `json:"owner_name"`
	MeterAddress    string  `json:"device_address"`
	PropertyID      uint64  `json:"property_id"`
	CadastralNumber string  `json:"cadastral_number"`
	City            string  `json:"city"`
	Street          string  `json:"street"`
	Number          string  `json:"number"`
}
