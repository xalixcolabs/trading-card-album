package card_pool_application_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"com.xalixcolabs.trading-card-album/context/card_pool/application"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestAssignCardMarksCardAsDrawn(t *testing.T) {
	mock := &queriermock.Querier{
		GetRandomAvailableCardFn: func(ctx context.Context, albumID string) (string, error) {
			return "card-5", nil
		},
	}

	drawn := false
	mock.MarkCardAsDrawnFn = func(ctx context.Context, arg sqlc.MarkCardAsDrawnParams) (sqlc.CardPool, error) {
		drawn = true
		return sqlc.CardPool{AlbumID: arg.AlbumID, CardID: arg.CardID, IsDrawn: 1}, nil
	}

	cardId, err := card_pool_application.AssignCard(mock, "album-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cardId != "card-5" {
		t.Errorf("expected card-5, got %s", cardId)
	}
	if !drawn {
		t.Error("expected card to be marked as drawn")
	}
}

func TestAssignCardResetsPoolWhenExhausted(t *testing.T) {
	calls := 0
	mock := &queriermock.Querier{
		GetRandomAvailableCardFn: func(ctx context.Context, albumID string) (string, error) {
			calls++
			if calls == 1 {
				return "", sql.ErrNoRows
			}
			return "card-2", nil
		},
	}

	reset := false
	mock.ResetCardPoolFn = func(ctx context.Context, albumID string) (sqlc.CardPool, error) {
		reset = true
		return sqlc.CardPool{}, nil
	}

	cardId, err := card_pool_application.AssignCard(mock, "album-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cardId != "card-2" {
		t.Errorf("expected card-2, got %s", cardId)
	}
	if !reset {
		t.Error("expected pool to be reset")
	}
}

func TestAssignCardReturnsErrorWhenResetFails(t *testing.T) {
	mock := &queriermock.Querier{
		GetRandomAvailableCardFn: func(ctx context.Context, albumID string) (string, error) {
			return "", sql.ErrNoRows
		},
		ResetCardPoolFn: func(ctx context.Context, albumID string) (sqlc.CardPool, error) {
			return sqlc.CardPool{}, errors.New("reset error")
		},
	}
	_, err := card_pool_application.AssignCard(mock, "album-1")
	if err == nil || err.Error() != "error al reiniciar el pool: reset error" {
		t.Errorf("expected reset error, got %v", err)
	}
}

func TestAssignCardReturnsErrorWhenLookupFails(t *testing.T) {
	mock := &queriermock.Querier{
		GetRandomAvailableCardFn: func(ctx context.Context, albumID string) (string, error) {
			return "", errors.New("lookup error")
		},
	}
	_, err := card_pool_application.AssignCard(mock, "album-1")
	if err == nil || err.Error() != "error al buscar tarjeta: lookup error" {
		t.Errorf("expected lookup error, got %v", err)
	}
}

func TestAssignCardReturnsErrorWhenMarkFails(t *testing.T) {
	mock := &queriermock.Querier{
		GetRandomAvailableCardFn: func(ctx context.Context, albumID string) (string, error) {
			return "card-5", nil
		},
		MarkCardAsDrawnFn: func(ctx context.Context, arg sqlc.MarkCardAsDrawnParams) (sqlc.CardPool, error) {
			return sqlc.CardPool{}, errors.New("mark error")
		},
	}
	_, err := card_pool_application.AssignCard(mock, "album-1")
	if err == nil || err.Error() != "error al marcar tarjeta como asignada: mark error" {
		t.Errorf("expected mark error, got %v", err)
	}
}