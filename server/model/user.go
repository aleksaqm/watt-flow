package model

type AccountStatus int

const (
	Active AccountStatus = iota
	Inactive
	Suspended
)

type User struct {
	Id           uint64        `gorm:"primary_key" json:"id"`
	Username     string        `gorm:"unique" json:"username"`
	Password     string        `json:"password,omitempty"`
	Email        string        `gorm:"unique" json:"email"`
	ProfileImage string        `json:"profile_image,omitempty"`
	Status       AccountStatus `json:"status"`
	Role         Role          `json:"role"`
	Meetings     []Meeting     `gorm:"foreignKey:ClerkID;foreignKey:UserID" json:"meetings,omitempty"` //maybe will need to separate clarkMeetings and userMeetings
}
