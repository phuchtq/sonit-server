package entity

import "time"

type User struct {
	UserId          string    `json:"user_id"`
	RoleId          string    `json:"role_id"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email" validate:"email, required"`
	Password        string    `json:"password" validate:"min=10, required"`
	ProfileAvatar   string    `json:"profile_avatar"`
	Gender          string    `json:"gender"`
	IsVip           bool      `json:"is_vip"`
	VipCode         *string   `json:"vip_code"`
	IsActive        bool      `json:"is_active"`
	IsActivated     bool      `json:"is_activated"`
	IsHaveToResetPw *bool     `json:"is_have_to_reset_password"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func GetUserTable() string {
	return "users"
}
