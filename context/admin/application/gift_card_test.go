package admin_application_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"com.xalixcolabs.trading-card-album/context/admin/application"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestGiftCardToUserCollectsCardWithoutChangingParticipant(t *testing.T) {
	var collected sqlc.CollectCardParams
	assignTouched := false
	markTouched := false

	mock := &queriermock.Querier{
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "u1"}, nil
		},
		GetCardFn: func(ctx context.Context, id string) (sqlc.Card, error) {
			return testAssignCard(), nil
		},
		CollectCardFn: func(ctx context.Context, arg sqlc.CollectCardParams) error {
			collected = arg
			return nil
		},
		UpdateAlbumParticipantAssignedCardFn: func(ctx context.Context, arg sqlc.UpdateAlbumParticipantAssignedCardParams) (sqlc.AlbumParticipant, error) {
			assignTouched = true
			return sqlc.AlbumParticipant{}, nil
		},
		MarkCardAsDrawnFn: func(ctx context.Context, arg sqlc.MarkCardAsDrawnParams) (sqlc.CardPool, error) {
			markTouched = true
			return sqlc.CardPool{}, nil
		},
	}

	card, err := admin_application.GiftCardToUser(mock, "u1", "c9")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if collected.UserID != "u1" || collected.AlbumID != "a1" || collected.CardID != "c9" {
		t.Errorf("unexpected collect params: %+v", collected)
	}
	if assignTouched {
		t.Error("expected the assigned card to remain unchanged")
	}
	if markTouched {
		t.Error("expected the gifted card NOT to be marked as drawn")
	}
	if card.ID != "c9" {
		t.Errorf("expected card c9, got %s", card.ID)
	}
}

func TestGiftCardToUserReturnsErrorWhenUserNotFound(t *testing.T) {
	mock := &queriermock.Querier{
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{}, sql.ErrNoRows
		},
	}
	_, err := admin_application.GiftCardToUser(mock, "u1", "c9")
	if err == nil || err.Error() != "Usuario no encontrado" {
		t.Errorf("expected 'Usuario no encontrado', got %v", err)
	}
}

func TestGiftCardToUserReturnsErrorWhenCardNotFound(t *testing.T) {
	mock := &queriermock.Querier{
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "u1"}, nil
		},
		GetCardFn: func(ctx context.Context, id string) (sqlc.Card, error) {
			return sqlc.Card{}, sql.ErrNoRows
		},
	}
	_, err := admin_application.GiftCardToUser(mock, "u1", "c9")
	if err == nil || err.Error() != "Tarjeta no encontrada" {
		t.Errorf("expected 'Tarjeta no encontrada', got %v", err)
	}
}

func TestGiftCardToUserReturnsErrorWhenCollectFails(t *testing.T) {
	mock := &queriermock.Querier{
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "u1"}, nil
		},
		GetCardFn: func(ctx context.Context, id string) (sqlc.Card, error) {
			return testAssignCard(), nil
		},
		CollectCardFn: func(ctx context.Context, arg sqlc.CollectCardParams) error {
			return errors.New("db error")
		},
	}
	_, err := admin_application.GiftCardToUser(mock, "u1", "c9")
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}
