package model

type City struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	CityName string `gorm:"type:varchar(100);not null" json:"city_name"`
}
