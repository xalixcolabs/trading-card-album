package admin_application

import (
	"context"
	"time"

	"com.xalixcolabs.trading-card-album/context/admin/model/dto"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func ListUsers(q database.Querier) ([]admin_dto.User, error) {
	ctx := context.Background()
	users, err := q.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]admin_dto.User, 0, len(users))
	for _, user := range users {
		result = append(result, mapUser(user))
	}
	return result, nil
}

func mapUser(user sqlc.User) admin_dto.User {
	return admin_dto.User{
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

func UpdateUserRole(q database.Querier, userId string, isAdmin bool) (admin_dto.User, error) {
	ctx := context.Background()
	isAdminInt := int64(0)
	if isAdmin {
		isAdminInt = 1
	}
	user, err := q.UpdateUserIsAdmin(ctx, sqlc.UpdateUserIsAdminParams{
		IsAdmin:   isAdminInt,
		UpdatedAt: time.Now().Unix(),
		ID:        userId,
	})
	if err != nil {
		return admin_dto.User{}, err
	}
	return mapUser(user), nil
}