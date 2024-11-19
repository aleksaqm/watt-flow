package model

type DeviceStatus struct {
	DeviceId string `gorm:"primaryKey" json:"device_id"`
	IsActive bool   `json:"is_active"`
}

func (DeviceStatus) TableName() string {
	return "device_status"
}
