package admin_application_test

import (
	"context"
	"errors"
	"testing"

	"com.xalixcolabs.trading-card-album/context/admin/application"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestGetOverview(t *testing.T) {
	mock := &queriermock.Querier{
		GetOverviewStatsFn: func(ctx context.Context) (sqlc.GetOverviewStatsRow, error) {
			return sqlc.GetOverviewStatsRow{
				Albums: 3, Users: 10, Cards: 50, Participants: 12, Contacts: 8, Collected: 30,
			}, nil
		},
	}
	overview, err := admin_application.GetOverview(mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if overview.Albums != 3 || overview.Users != 10 || overview.Cards != 50 ||
		overview.Participants != 12 || overview.Contacts != 8 || overview.Collected != 30 {
		t.Errorf("unexpected overview: %+v", overview)
	}
}

func TestGetOverviewReturnsError(t *testing.T) {
	mock := &queriermock.Querier{
		GetOverviewStatsFn: func(ctx context.Context) (sqlc.GetOverviewStatsRow, error) {
			return sqlc.GetOverviewStatsRow{}, errors.New("db error")
		},
	}
	if _, err := admin_application.GetOverview(mock); err == nil {
		t.Error("expected error")
	}
}

func TestListAlbums(t *testing.T) {
	mock := &queriermock.Querier{
		ListAlbumsWithStatsFn: func(ctx context.Context) ([]sqlc.ListAlbumsWithStatsRow, error) {
			return []sqlc.ListAlbumsWithStatsRow{
				{ID: "a1", Title: "DevFest", TotalCards: 20, CreatedAt: 1, ParticipantCount: 5},
			}, nil
		},
	}
	albums, err := admin_application.ListAlbums(mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(albums) != 1 || albums[0].ID != "a1" || albums[0].ParticipantCount != 5 {
		t.Errorf("unexpected albums: %+v", albums)
	}
}

func TestDeleteAlbum(t *testing.T) {
	var deletedID string
	mock := &queriermock.Querier{
		DeleteAlbumFn: func(ctx context.Context, id string) error {
			deletedID = id
			return nil
		},
	}
	if err := admin_application.DeleteAlbum(mock, "a1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deletedID != "a1" {
		t.Errorf("expected deleted id a1, got %s", deletedID)
	}
}

func TestListUsers(t *testing.T) {
	mock := &queriermock.Querier{
		ListUsersFn: func(ctx context.Context) ([]sqlc.User, error) {
			return []sqlc.User{
				{ID: "u1", Email: "a@b.c", IsAdmin: 1},
				{ID: "u2", Email: "d@e.f", IsAdmin: 0},
			}, nil
		},
	}
	users, err := admin_application.ListUsers(mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 || users[0].IsAdmin != 1 || users[1].IsAdmin != 0 {
		t.Errorf("unexpected users: %+v", users)
	}
}

func TestUpdateUserRole(t *testing.T) {
	mock := &queriermock.Querier{
		UpdateUserIsAdminFn: func(ctx context.Context, arg sqlc.UpdateUserIsAdminParams) (sqlc.User, error) {
			return sqlc.User{ID: arg.ID, IsAdmin: arg.IsAdmin}, nil
		},
	}
	user, err := admin_application.UpdateUserRole(mock, "u1", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "u1" || user.IsAdmin != 1 {
		t.Errorf("expected admin role, got %+v", user)
	}
}

func TestListCards(t *testing.T) {
	mock := &queriermock.Querier{
		ListCardsFn: func(ctx context.Context) ([]sqlc.Card, error) {
			return []sqlc.Card{{ID: "c1", Name: "Gopher", Number: "01"}}, nil
		},
	}
	cards, err := admin_application.ListCards(mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cards) != 1 || cards[0].ID != "c1" {
		t.Errorf("unexpected cards: %+v", cards)
	}
}

func TestDeleteCard(t *testing.T) {
	var deletedID string
	mock := &queriermock.Querier{
		DeleteCardFn: func(ctx context.Context, id string) error {
			deletedID = id
			return nil
		},
	}
	if err := admin_application.DeleteCard(mock, "c1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deletedID != "c1" {
		t.Errorf("expected deleted id c1, got %s", deletedID)
	}
}