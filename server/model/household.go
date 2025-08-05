package model

import "errors"

type HouseholdStatus int

const (
	InactiveHousehold HouseholdStatus = iota
	OwnedHousehold
	UnownedHousehold
)

func (r HouseholdStatus) HouseholdStatusToString() string {
	switch r {
	case 0:
		return "Inactive"
	case 1:
		return "Owned"
	case 2:
		return "Unowned"
	default:
		return "Unknown"
	}
}

func ParseHouseholdStatus(status string) (HouseholdStatus, error) {
	switch status {
	case "InactiveHousehold":
		return InactiveHousehold, nil
	case "OwnedHousehold":
		return OwnedHousehold, nil
	case "UnownedHousehold":
		return UnownedHousehold, nil
	default:
		return InactiveHousehold, errors.New("invalid status value")
	}
}

type Household struct {
	Id              uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Floor           uint32          `json:"floor"`
	Suite           string          `json:"suite"`
	Status          HouseholdStatus `json:"status"`
	SqFootage       float32         `json:"sq_footage"`
	OwnerID         *uint64         `gorm:"column:owner_id;null" json:"owner_id"`
	Owner           *User           `gorm:"foreignKey:OwnerID" json:"owner"`
	DeviceStatus    *DeviceStatus   `gorm:"foreignKey:DeviceStatusID" json:"device_status"`
	DeviceStatusID  string          `gorm:"column:device_status_id" json:"device_status_id"`
	PropertyID      uint64          `gorm:"column:property_id" json:"property_id"`
	CadastralNumber string          `gorm:"unique" json:"cadastral_number"`
	Property        Property        `gorm:"foreignKey:PropertyID" json:"property"`
}
