package model

type DeviceStatus struct {
	Address     string `gorm:"primaryKey" json:"address"`
	IsActive    bool   `json:"is_active"`
	HouseholdID uint64 `gorm:"column:household_id;unique" json:"household_id"`
}
