package album_application_test

import (
	"context"
	"errors"
	"testing"

	"com.xalixcolabs.trading-card-album/context/album/application"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestGetAlbumsByUserReturnsAlbums(t *testing.T) {
	mock := &queriermock.Querier{
		ListAlbumsByUserIdFn: func(ctx context.Context, userID string) ([]sqlc.Album, error) {
			if userID != "user-1" {
				t.Errorf("expected user id user-1, got %s", userID)
			}
			return []sqlc.Album{
				{ID: "album-1", Title: "Album 1", TotalCards: 2, CreatedAt: 100},
				{ID: "album-2", Title: "Album 2", TotalCards: 5, CreatedAt: 200},
			}, nil
		},
	}

	albums, err := album_application.GetAlbumsByUser(mock, testUser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(albums) != 2 {
		t.Fatalf("expected 2 albums, got %d", len(albums))
	}
	if albums[0].ID != "album-1" || albums[1].ID != "album-2" {
		t.Errorf("unexpected albums order: %+v", albums)
	}
}

func TestGetAlbumsByUserReturnsError(t *testing.T) {
	mock := &queriermock.Querier{
		ListAlbumsByUserIdFn: func(ctx context.Context, userID string) ([]sqlc.Album, error) {
			return nil, errors.New("db error")
		},
	}
	_, err := album_application.GetAlbumsByUser(mock, testUser())
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}