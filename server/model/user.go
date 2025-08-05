package model

import "fmt"

type AccountStatus int

const (
	Active AccountStatus = iota
	Inactive
	Suspended
)

func (s AccountStatus) StatusToString() string {
	switch s {
	case 0:
		return "Active"
	case 1:
		return "Inactive"
	case 2:
		return "Suspended"
	default:
		return "Unknown"
	}
}

func ParseAccountStatus(status string) (Role, error) {
	switch status {
	case "Active":
		return 0, nil
	case "Inactive":
		return 1, nil
	case "Suspended":
		return 2, nil
	default:
		return -1, fmt.Errorf("invalid account status: %s", status)
	}
}

type User struct {
	Id           uint64        `gorm:"primary_key" json:"id"`
	Username     string        `gorm:"unique" json:"username"`
	FirstName    string        `json:"first_name"`
	LastName     string        `json:"last_name"`
	Password     string        `json:"password,omitempty"`
	Email        string        `json:"email"`
	ProfileImage string        `json:"profile_image,omitempty"`
	Status       AccountStatus `json:"status"`
	Role         Role          `json:"role"`
	Meetings     []Meeting     `gorm:"foreignKey:ClerkID;foreignKey:UserID" json:"meetings,omitempty"` // maybe will need to separate clarkMeetings and userMeetings
}
