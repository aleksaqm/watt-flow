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
		return 0, nil
	case "Clerk":
		return 1, nil
	case "Admin":
		return 2, nil
	case "SuperAdmin":
		return 3, nil
	default:
		return -1, fmt.Errorf("invalid role: %s", roleStr)
	}
}
