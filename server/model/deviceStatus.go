package model

type DeviceStatus struct {
	Address     string `gorm:"primaryKey" json:"address"`
	IsActive    bool   `json:"is_active"`
	HouseholdID uint64 `gorm:"column:household_id" json:"household_id"`
}

func (DeviceStatus) TableName() string {
	return "device_status"
}
