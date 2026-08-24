package user_application

import (
	"context"

	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database"
)

func GetUserById(q database.Querier, id string) (user_model.User, error) {
	ctx := context.Background()
	user, err := q.GetUser(ctx, id)
	return user_model.NewUserFromSqlcUser(user), err
}
