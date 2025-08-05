package model

type Address struct {
	Id        uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	City      string  `json:"city"`
	Street    string  `json:"street"`
	Number    string  `json:"number"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
