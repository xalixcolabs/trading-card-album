package album_application_test

import (
	"context"
	"errors"
	"testing"

	"com.xalixcolabs.trading-card-album/context/album/application"
	dto "com.xalixcolabs.trading-card-album/context/album/model/dto"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestCreateAlbumCreatesCardsAndPoolRows(t *testing.T) {
	mock := &queriermock.Querier{
		CreateAlbumFn: func(ctx context.Context, arg sqlc.CreateAlbumParams) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1", Title: arg.Title, TotalCards: arg.TotalCards, CreatedAt: arg.CreatedAt}, nil
		},
		CreateCardFn: func(ctx context.Context, arg sqlc.CreateCardParams) (sqlc.Card, error) {
			return sqlc.Card{
				ID:          "card-" + arg.Number,
				AlbumID:     arg.AlbumID,
				Number:      arg.Number,
				Name:        arg.Name,
				Description: arg.Description,
				ImageUrl:    arg.ImageUrl,
				CreatedAt:   arg.CreatedAt,
				UpdatedAt:   arg.UpdatedAt,
			}, nil
		},
		CreateCardPoolRowFn: func(ctx context.Context, arg sqlc.CreateCardPoolRowParams) (sqlc.CardPool, error) {
			return sqlc.CardPool{AlbumID: arg.AlbumID, CardID: arg.CardID}, nil
		},
	}

	request := dto.CreateAlbumRequest{
		Title: "DevFest 2026",
		Cards: []dto.CreateAlbumCardRequest{
			{Number: "01", Name: "Card A", Description: "desc A", ImageUrl: "http://localhost:8081/a.webp"},
			{Number: "02", Name: "Card B", Description: "desc B", ImageUrl: "http://localhost:8081/b.webp"},
		},
	}

	album, err := album_application.CreateAlbum(mock, request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if album.ID != "album-1" {
		t.Errorf("expected album id album-1, got %s", album.ID)
	}
	if album.Title != "DevFest 2026" {
		t.Errorf("expected title DevFest 2026, got %s", album.Title)
	}
	if album.TotalCards != 2 {
		t.Errorf("expected total cards 2, got %d", album.TotalCards)
	}
	if len(album.Cards) != 2 {
		t.Fatalf("expected 2 cards, got %d", len(album.Cards))
	}
}

func TestCreateAlbumReturnsErrorWhenCreateCardFails(t *testing.T) {
	mock := &queriermock.Querier{
		CreateAlbumFn: func(ctx context.Context, arg sqlc.CreateAlbumParams) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1", Title: arg.Title}, nil
		},
		CreateCardFn: func(ctx context.Context, arg sqlc.CreateCardParams) (sqlc.Card, error) {
			return sqlc.Card{}, errors.New("create card error")
		},
	}

	request := dto.CreateAlbumRequest{
		Title: "DevFest 2026",
		Cards: []dto.CreateAlbumCardRequest{{Number: "01", Name: "Card A"}},
	}

	_, err := album_application.CreateAlbum(mock, request)
	if err == nil || err.Error() != "create card error" {
		t.Errorf("expected create card error, got %v", err)
	}
}

func TestCreateAlbumReturnsErrorWhenCreatePoolRowFails(t *testing.T) {
	mock := &queriermock.Querier{
		CreateAlbumFn: func(ctx context.Context, arg sqlc.CreateAlbumParams) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1", Title: arg.Title}, nil
		},
		CreateCardFn: func(ctx context.Context, arg sqlc.CreateCardParams) (sqlc.Card, error) {
			return sqlc.Card{ID: "card-1", Number: arg.Number}, nil
		},
		CreateCardPoolRowFn: func(ctx context.Context, arg sqlc.CreateCardPoolRowParams) (sqlc.CardPool, error) {
			return sqlc.CardPool{}, errors.New("pool error")
		},
	}

	request := dto.CreateAlbumRequest{
		Title: "DevFest 2026",
		Cards: []dto.CreateAlbumCardRequest{{Number: "01", Name: "Card A"}},
	}

	_, err := album_application.CreateAlbum(mock, request)
	if err == nil || err.Error() != "pool error" {
		t.Errorf("expected pool error, got %v", err)
	}
}