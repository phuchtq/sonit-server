package request

type UpdateRoleRequest struct {
	RoleId   string `json:"role_id" validate:"required"`
	RoleName string `json:"role_name"`
}
