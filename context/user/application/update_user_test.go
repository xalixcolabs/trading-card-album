package user_application_test

import (
	"context"
	"errors"
	"testing"

	"com.xalixcolabs.trading-card-album/context/user/application"
	dto "com.xalixcolabs.trading-card-album/context/user/model/dto"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestUpdateUserAppliesRequest(t *testing.T) {
	var captured sqlc.UpdateUserParams
	mock := &queriermock.Querier{
		UpdateUserFn: func(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
			captured = arg
			return sqlc.User{
				ID:          arg.ID,
				Name:        arg.Name,
				Email:       arg.Email,
				Github:      arg.Github,
				Linkedin:    arg.Linkedin,
				Web:         arg.Web,
				Description: arg.Description,
				UpdatedAt:   arg.UpdatedAt,
			}, nil
		},
	}

	request := dto.UpdateUserRequest{
		Name:        "Uziel",
		Email:       "uziel@example.com",
		Github:      "github.com/uziel",
		Linkedin:    "linkedin.com/uziel",
		Web:         "uziel.dev",
		Description: "Go dev",
	}

	user, err := user_application.UpdateUser(mock, "user-1", request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.ID != "user-1" {
		t.Errorf("expected id user-1, got %s", captured.ID)
	}
	if user.Name != "Uziel" || user.Description != "Go dev" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestUpdateUserReturnsError(t *testing.T) {
	mock := &queriermock.Querier{
		UpdateUserFn: func(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
			return sqlc.User{}, errors.New("db error")
		},
	}
	_, err := user_application.UpdateUser(mock, "user-1", dto.UpdateUserRequest{})
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}