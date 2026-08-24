package contact_dto

type Contact struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Github      string `json:"github"`
	Linkedin    string `json:"linkedin"`
	Web         string `json:"web"`
	Description string `json:"description"`
	Picture     string `json:"picture"`
	ScannedAt   int64  `json:"scanned_at"`
}