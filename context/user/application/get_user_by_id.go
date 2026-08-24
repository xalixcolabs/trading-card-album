package user_application

import (
	"context"

	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func GetUserById(id string) (user_model.User, error) {
	ctx := context.Background()
	queries := sqlc.New(database.GetDatabase())
	user, err := queries.GetUser(ctx, id)
	return user_model.NewUserFromSqlcUser(user), err
}
