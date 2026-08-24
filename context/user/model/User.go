package user_model

import "com.xalixcolabs.trading-card-album/database/sqlc"

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Github      string `json:"github"`
	Linkedin    string `json:"linkedin"`
	Web         string `json:"web"`
	Description string `json:"description"`
	IsAdmin     int64  `json:"-"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"update_at"`
}

func NewUserFromSqlcUser(user sqlc.User) User {
	return User{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Github:      user.Github,
		Linkedin:    user.Linkedin,
		Web:         user.Web,
		Description: user.Description,
		IsAdmin:     user.IsAdmin,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
