package user_application_test

import (
	"context"
	"errors"
	"testing"

	"com.xalixcolabs.trading-card-album/context/user/application"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestCreateUserSetsAdminFlag(t *testing.T) {
	var captured sqlc.CreateUserParams
	mock := &queriermock.Querier{
		CreateUserFn: func(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
			captured = arg
			return sqlc.User{
				ID:      arg.ID,
				Email:   arg.Email,
				IsAdmin: arg.IsAdmin,
			}, nil
		},
	}

	user, err := user_application.CreateUser(mock, "admin@example.com", true, "https://pic.example/1.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "admin@example.com" {
		t.Errorf("expected email admin@example.com, got %s", user.Email)
	}
	if user.IsAdmin != 1 {
		t.Errorf("expected is_admin 1, got %d", user.IsAdmin)
	}
	if captured.IsAdmin != 1 {
		t.Errorf("expected captured is_admin 1, got %d", captured.IsAdmin)
	}
}

func TestCreateUserNonAdminFlag(t *testing.T) {
	mock := &queriermock.Querier{
		CreateUserFn: func(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
			return sqlc.User{ID: arg.ID, Email: arg.Email, IsAdmin: arg.IsAdmin}, nil
		},
	}
	user, err := user_application.CreateUser(mock, "user@example.com", false, "https://pic.example/2.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.IsAdmin != 0 {
		t.Errorf("expected is_admin 0, got %d", user.IsAdmin)
	}
}

func TestCreateUserReturnsError(t *testing.T) {
	mock := &queriermock.Querier{
		CreateUserFn: func(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
			return sqlc.User{}, errors.New("db error")
		},
	}
	_, err := user_application.CreateUser(mock, "user@example.com", false, "https://pic.example/2.png")
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}