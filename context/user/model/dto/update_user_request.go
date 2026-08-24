package user_dto

type UpdateUserRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Github      string `json:"github"`
	Linkedin    string `json:"linkedin"`
	Web         string `json:"web"`
	Description string `json:"description"`
}
