package admin_application

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"com.xalixcolabs.trading-card-album/context/admin/model/dto"
	card_model "com.xalixcolabs.trading-card-album/context/card/model"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func ListUsers(q database.Querier) ([]admin_dto.User, error) {
	ctx := context.Background()
	users, err := q.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return mapUsers(users), nil
}

func SearchUsers(q database.Querier, email string) ([]admin_dto.User, error) {
	ctx := context.Background()
	users, err := q.SearchUsersByEmail(ctx, "%"+email+"%")
	if err != nil {
		return nil, err
	}
	return mapUsers(users), nil
}

func mapUsers(users []sqlc.User) []admin_dto.User {
	result := make([]admin_dto.User, 0, len(users))
	for _, user := range users {
		result = append(result, mapUser(user))
	}
	return result
}

func GetUserDetail(q database.Querier, userId string) (admin_dto.UserDetail, error) {
	ctx := context.Background()
	user, err := q.GetUser(ctx, userId)
	if errors.Is(err, sql.ErrNoRows) {
		return admin_dto.UserDetail{}, err
	}
	if err != nil {
		return admin_dto.UserDetail{}, err
	}

	albums, err := q.ListAlbumsByUserId(ctx, userId)
	if err != nil {
		return admin_dto.UserDetail{}, err
	}
	albumDTOs := make([]admin_dto.Album, 0, len(albums))
	for _, album := range albums {
		albumDTOs = append(albumDTOs, admin_dto.Album{
			ID:         album.ID,
			Title:      album.Title,
			TotalCards: album.TotalCards,
			CreatedAt:  album.CreatedAt,
		})
	}

	cards, err := q.ListCollectedCardsByUser(ctx, userId)
	if err != nil {
		return admin_dto.UserDetail{}, err
	}

	return admin_dto.UserDetail{
		User:   mapUser(user),
		Albums: albumDTOs,
		Cards:  card_model.NewCardSliceFromSqlcCardSlice(cards),
	}, nil
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
		Picture:     user.Picture,
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