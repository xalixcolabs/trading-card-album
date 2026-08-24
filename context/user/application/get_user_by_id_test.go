package user_application_test

import (
	"context"
	"database/sql"
	"testing"

	"com.xalixcolabs.trading-card-album/context/user/application"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestGetUserByIdReturnsUser(t *testing.T) {
	mock := &queriermock.Querier{
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			if id != "user-1" {
				t.Errorf("expected id user-1, got %s", id)
			}
			return sqlc.User{ID: "user-1", Email: "user@example.com", Name: "Uziel"}, nil
		},
	}

	user, err := user_application.GetUserById(mock, "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "user-1" || user.Email != "user@example.com" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestGetUserByIdReturnsError(t *testing.T) {
	mock := &queriermock.Querier{
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{}, sql.ErrNoRows
		},
	}
	_, err := user_application.GetUserById(mock, "missing")
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}