package album_application_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"com.xalixcolabs.trading-card-album/context/album/application"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func testParticipantUser() user_model.User {
	return user_model.User{ID: "user-1"}
}

func TestGetAlbumByIdReturnsOnlyCollectedCards(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			if id != "album-1" {
				t.Errorf("expected id album-1, got %s", id)
			}
			return sqlc.Album{ID: "album-1", Title: "DevFest 2026", TotalCards: 3, CreatedAt: 100}, nil
		},
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			if arg.AlbumID != "album-1" || arg.UserID != "user-1" {
				t.Errorf("unexpected params: %+v", arg)
			}
			return sqlc.AlbumParticipant{AlbumID: "album-1", UserID: "user-1", AssignedCardID: "card-1"}, nil
		},
		GetUserCollectionFn: func(ctx context.Context, arg sqlc.GetUserCollectionParams) ([]sqlc.Card, error) {
			if arg.UserID != "user-1" || arg.AlbumID != "album-1" {
				t.Errorf("unexpected params: %+v", arg)
			}
			return []sqlc.Card{
				{ID: "card-1", AlbumID: "album-1", Number: "01", Name: "Gopher"},
			}, nil
		},
	}

	album, err := album_application.GetAlbumById(mock, testParticipantUser(), "album-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if album.ID != "album-1" {
		t.Errorf("expected id album-1, got %s", album.ID)
	}
	if album.Title != "DevFest 2026" {
		t.Errorf("expected title DevFest 2026, got %s", album.Title)
	}
	if album.TotalCards != 3 {
		t.Errorf("expected total cards 3, got %d", album.TotalCards)
	}
	if len(album.Cards) != 1 {
		t.Fatalf("expected only collected card, got %d", len(album.Cards))
	}
	if album.Cards[0].ID != "card-1" {
		t.Errorf("expected collected card card-1, got %+v", album.Cards[0])
	}
}

func TestGetAlbumByIdReturnsErrorWhenAlbumNotFound(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{}, sql.ErrNoRows
		},
	}
	_, err := album_application.GetAlbumById(mock, testParticipantUser(), "missing")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestGetAlbumByIdReturnsErrorWhenNotParticipant(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1", Title: "Album"}, nil
		},
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return sqlc.AlbumParticipant{}, sql.ErrNoRows
		},
	}
	_, err := album_application.GetAlbumById(mock, testParticipantUser(), "album-1")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestGetAlbumByIdReturnsErrorOnLookupFailure(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{}, errors.New("db error")
		},
	}
	_, err := album_application.GetAlbumById(mock, testParticipantUser(), "album-1")
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestGetAlbumByIdReturnsErrorWhenCollectionFails(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1", Title: "Album"}, nil
		},
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return sqlc.AlbumParticipant{AlbumID: "album-1", UserID: "user-1"}, nil
		},
		GetUserCollectionFn: func(ctx context.Context, arg sqlc.GetUserCollectionParams) ([]sqlc.Card, error) {
			return nil, errors.New("collection error")
		},
	}
	_, err := album_application.GetAlbumById(mock, testParticipantUser(), "album-1")
	if err == nil || err.Error() != "collection error" {
		t.Errorf("expected collection error, got %v", err)
	}
}