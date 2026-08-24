package admin_dto

type UpdateUserRoleRequest struct {
	IsAdmin bool `json:"is_admin"`
}