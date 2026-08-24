package card_application_test

import (
	"context"
	"database/sql"
	"testing"

	"com.xalixcolabs.trading-card-album/context/card/application"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestGetCardByAlbumIdReturnsAssignedCard(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return sqlc.AlbumParticipant{AlbumID: "album-1", UserID: "user-1", AssignedCardID: "card-3"}, nil
		},
		GetCardByAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetCardByAlbumParticipantParams) (sqlc.Card, error) {
			return sqlc.Card{ID: "card-3", AlbumID: "album-1", Number: "03", Name: "Flame"}, nil
		},
	}

	card, err := card_application.GetCardByAlbumId(mock, user_model.User{ID: "user-1"}, "album-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if card.ID != "card-3" || card.Name != "Flame" {
		t.Errorf("unexpected card: %+v", card)
	}
}

func TestGetCardByAlbumIdReturnsEmptyCardWhenNotParticipant(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return sqlc.AlbumParticipant{}, sql.ErrNoRows
		},
	}
	card, err := card_application.GetCardByAlbumId(mock, user_model.User{ID: "user-1"}, "album-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if card.ID != "" {
		t.Errorf("expected empty card, got %+v", card)
	}
}

func TestGetCardByAlbumIdReturnsEmptyCardWhenCardMissing(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return sqlc.AlbumParticipant{AlbumID: "album-1", UserID: "user-1", AssignedCardID: "card-3"}, nil
		},
		GetCardByAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetCardByAlbumParticipantParams) (sqlc.Card, error) {
			return sqlc.Card{}, sql.ErrNoRows
		},
	}
	card, err := card_application.GetCardByAlbumId(mock, user_model.User{ID: "user-1"}, "album-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if card.ID != "" {
		t.Errorf("expected empty card, got %+v", card)
	}
}