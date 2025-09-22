package entity

import "time"

type Role struct {
	RoleId       string    `json:"role_id"`
	RoleName     string    `json:"role_name"`
	ActiveStatus bool      `json:"active_status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func GetRoleTable() string {
	return "roles"
}
