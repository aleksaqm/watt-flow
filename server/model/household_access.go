package model

import (
	"time"
)

type HouseholdAccess struct {
	Id          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	HouseholdID uint64    `gorm:"uniqueIndex:idx_household_user" json:"householdId"`
	UserID      uint64    `gorm:"uniqueIndex:idx_household_user" json:"userId"`
	GrantedAt   time.Time `gorm:"autoCreateTime" json:"grantedAt"`

	User      User      `gorm:"foreignKey:UserID"`
	Household Household `gorm:"foreignKey:HouseholdID"`
}
