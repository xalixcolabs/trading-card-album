package card_pool_application_test

import (
	"context"
	"errors"
	"testing"

	"com.xalixcolabs.trading-card-album/context/card_pool/application"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestGetMyCardsByAlbumIdReturnsCollection(t *testing.T) {
	mock := &queriermock.Querier{
		GetUserCollectionFn: func(ctx context.Context, arg sqlc.GetUserCollectionParams) ([]sqlc.Card, error) {
			if arg.UserID != "user-1" || arg.AlbumID != "album-1" {
				t.Errorf("unexpected params: %+v", arg)
			}
			return []sqlc.Card{
				{ID: "card-1", Number: "01", Name: "Card One"},
				{ID: "card-2", Number: "02", Name: "Card Two"},
			}, nil
		},
	}

	cards, err := card_pool_application.GetMyCardsByAlbumId(mock, user_model.User{ID: "user-1"}, "album-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cards) != 2 {
		t.Fatalf("expected 2 cards, got %d", len(cards))
	}
	if cards[0].ID != "card-1" || cards[1].ID != "card-2" {
		t.Errorf("unexpected cards: %+v", cards)
	}
}

func TestGetMyCardsByAlbumIdReturnsEmptyOnError(t *testing.T) {
	mock := &queriermock.Querier{
		GetUserCollectionFn: func(ctx context.Context, arg sqlc.GetUserCollectionParams) ([]sqlc.Card, error) {
			return nil, errors.New("db error")
		},
	}
	cards, err := card_pool_application.GetMyCardsByAlbumId(mock, user_model.User{ID: "user-1"}, "album-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(cards) != 0 {
		t.Errorf("expected empty cards, got %+v", cards)
	}
}