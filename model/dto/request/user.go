package request

// Security request

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginSecurityRequest struct {
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ---------------------------------------

type CreateUserRequest struct {
	ActorId       string `json:"actor_id"`
	RoleId        string `json:"role_id"`
	FullName      string `json:"full_name" validate:"required"`
	Email         string `json:"email" validate:"required"`
	Password      string `json:"password" validate:"required,min=8"`
	ProfileAvatar string `json:"profile_avatar"`
	Gender        string `json:"gender"`
}

type UpdateUserRequest struct {
	ActorId       string `json:"actor_id" validate:"required"`
	UserId        string `json:"user_id" validate:"required"`
	RoleId        string `json:"role_id"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	Password      string `json:"password" validate:"min=10"`
	ProfileAvatar string `json:"profile_avatar"`
	Gender        string `json:"gender"`
	IsVip         *bool  `json:"is_vip"`
}

type ChangeUserStatusRequest struct {
	ActorId string `json:"actor_id" validate:"required"`
	UserId  string `json:"user_id" validate:"required"`
	Status  *bool  `json:"status"`
}

// For vip users

type CreateVipCodeRequest struct {
	UserId  string `json:"user_id" validate:"required"`
	VipCode string `json:"vip_code" validate:"required, max=4, min=4"`
}
