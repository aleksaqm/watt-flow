package dto

type SearchHouseholdDto struct {
	Street string `json:"street"`
	Number string `json:"number"`
	City   string `json:"city"`
	Id     string `json:"household_id"`
	Floor  uint32 `json:"floor"`
	Suite  string `json:"suite"`
}
