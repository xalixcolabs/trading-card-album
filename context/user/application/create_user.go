package user_application

import (
	"context"
	"time"

	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
	"github.com/matoous/go-nanoid/v2"
)

func CreateUser(email string, isAdmin bool) (sqlc.User, error) {
	ctx := context.Background()
	queries := sqlc.New(database.GetDatabase())
	id, _ := gonanoid.New()
	isAdmin_ := 0
	if isAdmin {
		isAdmin_ = 1
	}
	user, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:          id,
		Name:        "",
		Email:       email,
		Github:      "",
		Linkedin:    "",
		Web:         "",
		IsAdmin:     int64(isAdmin_),
		Description: "",

		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})
	return user, err
}
