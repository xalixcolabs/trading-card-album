package admin_dto

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Github      string `json:"github"`
	Linkedin    string `json:"linkedin"`
	Web         string `json:"web"`
	Description string `json:"description"`
	IsAdmin     int64  `json:"is_admin"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}