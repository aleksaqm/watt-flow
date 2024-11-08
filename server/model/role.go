package model

import "fmt"

type Role int

const (
	Regular Role = iota
	Clerk
	Admin
	SuperAdmin
)

func (r Role) RoleToString() string {
	switch r {
	case 0:
		return "Regular"
	case 1:
		return "Clerk"
	case 2:
		return "Admin"
	case 3:
		return "SuperAdmin"
	default:
		return "Unknown"
	}
}

func ParseRole(roleStr string) (Role, error) {
	switch roleStr {
	case "Regular":
		return Regular, nil
	case "Clerk":
		return Clerk, nil
	case "Admin":
		return Admin, nil
	case "SuperAdmin":
		return SuperAdmin, nil
	default:
		return -1, fmt.Errorf("invalid role: %s", roleStr)
	}
}
