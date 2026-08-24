package user_application

import (
	"context"
	"time"

	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/context/user/model/dto"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func UpdateUser(id string, request user_dto.UpdateUserRequest) (user_model.User, error) {
	ctx := context.Background()
	queries := sqlc.New(database.GetDatabase())
	user, err := queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:          id,
		Name:        request.Name,
		Email:       request.Email,
		Github:      request.Github,
		Linkedin:    request.Linkedin,
		Web:         request.Web,
		Description: request.Description,
		UpdatedAt:   time.Now().Unix(),
	})
	return user_model.NewUserFromSqlcUser(user), err
}
