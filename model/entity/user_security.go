package entity

import "time"

type UserSecurity struct {
	UserId       string     `json:"user_id"`
	AccessToken  *string    `json:"access_token"`
	RefreshToken *string    `json:"refresh_token"`
	ActionToken  *string    `json:"action_token"` // Lưu trữ token cho thay đổi mail, reset pass
	FailAccess   int        `json:"fail_access"`
	LastFail     *time.Time `json:"last_fail"`
}

func GetUserSecurityTable() string {
	return "user_securities"
}
