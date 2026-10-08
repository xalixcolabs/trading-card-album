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

func testAssignCard() sqlc.Card {
	return sqlc.Card{
		ID:          "c9",
		AlbumID:     "a1",
		Number:      "09",
		Name:        "Gopher",
		Description: "desc",
		ImageUrl:    "http://localhost:8081/c9.webp",
	}
}

func TestAssignCardToUserUpdatesExistingParticipant(t *testing.T) {
	var updated sqlc.UpdateAlbumParticipantAssignedCardParams
	marked := false
	collected := false

	mock := &queriermock.Querier{
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "u1"}, nil
		},
		GetCardFn: func(ctx context.Context, id string) (sqlc.Card, error) {
			return testAssignCard(), nil
		},
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return sqlc.AlbumParticipant{AlbumID: "a1", UserID: "u1", AssignedCardID: "c1", Secret: "s"}, nil
		},
		UpdateAlbumParticipantAssignedCardFn: func(ctx context.Context, arg sqlc.UpdateAlbumParticipantAssignedCardParams) (sqlc.AlbumParticipant, error) {
			updated = arg
			return sqlc.AlbumParticipant{AlbumID: arg.AlbumID, UserID: arg.UserID, AssignedCardID: arg.AssignedCardID, Secret: "s"}, nil
		},
		MarkCardAsDrawnFn: func(ctx context.Context, arg sqlc.MarkCardAsDrawnParams) (sqlc.CardPool, error) {
			marked = true
			return sqlc.CardPool{AlbumID: arg.AlbumID, CardID: arg.CardID, IsDrawn: 1}, nil
		},
	}
	mock.CollectCardFn = func(ctx context.Context, arg sqlc.CollectCardParams) error {
		collected = true
		return nil
	}

	card, err := admin_application.AssignCardToUser(mock, "u1", "c9")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.AssignedCardID != "c9" || updated.AlbumID != "a1" || updated.UserID != "u1" {
		t.Errorf("unexpected update params: %+v", updated)
	}
	if !marked {
		t.Error("expected card to be marked as drawn")
	}
	if !collected {
		t.Error("expected card to be collected")
	}
	if card.ID != "c9" {
		t.Errorf("expected card c9, got %s", card.ID)
	}
}

func TestAssignCardToUserCreatesParticipantWhenMissing(t *testing.T) {
	var created sqlc.CreateAlbumParticipantParams
	createCalled := false
	collected := false

	mock := &queriermock.Querier{
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "u1"}, nil
		},
		GetCardFn: func(ctx context.Context, id string) (sqlc.Card, error) {
			return testAssignCard(), nil
		},
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return sqlc.AlbumParticipant{}, sql.ErrNoRows
		},
		CreateAlbumParticipantFn: func(ctx context.Context, arg sqlc.CreateAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			created = arg
			createCalled = true
			return sqlc.AlbumParticipant{
				AlbumID:        arg.AlbumID,
				UserID:         arg.UserID,
				AssignedCardID: arg.AssignedCardID,
				JoinedAt:       arg.JoinedAt,
				Secret:         arg.Secret,
			}, nil
		},
		MarkCardAsDrawnFn: func(ctx context.Context, arg sqlc.MarkCardAsDrawnParams) (sqlc.CardPool, error) {
			return sqlc.CardPool{AlbumID: arg.AlbumID, CardID: arg.CardID, IsDrawn: 1}, nil
		},
	}
	mock.CollectCardFn = func(ctx context.Context, arg sqlc.CollectCardParams) error {
		collected = true
		return nil
	}

	card, err := admin_application.AssignCardToUser(mock, "u1", "c9")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !createCalled {
		t.Fatal("expected participant to be created")
	}
	if created.AssignedCardID != "c9" || created.AlbumID != "a1" || created.UserID != "u1" {
		t.Errorf("unexpected create params: %+v", created)
	}
	if created.Secret == "" {
		t.Error("expected a secret to be generated for the participant")
	}
	if !collected {
		t.Error("expected card to be collected")
	}
	if card.ID != "c9" {
		t.Errorf("expected card c9, got %s", card.ID)
	}
}

func TestAssignCardToUserReturnsErrorWhenUserNotFound(t *testing.T) {
	mock := &queriermock.Querier{
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{}, sql.ErrNoRows
		},
	}
	_, err := admin_application.AssignCardToUser(mock, "u1", "c9")
	if err == nil || err.Error() != "Usuario no encontrado" {
		t.Errorf("expected 'Usuario no encontrado', got %v", err)
	}
}

func TestAssignCardToUserReturnsErrorWhenCardNotFound(t *testing.T) {
	mock := &queriermock.Querier{
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "u1"}, nil
		},
		GetCardFn: func(ctx context.Context, id string) (sqlc.Card, error) {
			return sqlc.Card{}, sql.ErrNoRows
		},
	}
	_, err := admin_application.AssignCardToUser(mock, "u1", "c9")
	if err == nil || err.Error() != "Tarjeta no encontrada" {
		t.Errorf("expected 'Tarjeta no encontrada', got %v", err)
	}
}

func TestAssignCardToUserReturnsErrorWhenUpdateFails(t *testing.T) {
	mock := &queriermock.Querier{
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "u1"}, nil
		},
		GetCardFn: func(ctx context.Context, id string) (sqlc.Card, error) {
			return testAssignCard(), nil
		},
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return sqlc.AlbumParticipant{AlbumID: "a1", UserID: "u1", AssignedCardID: "c1"}, nil
		},
		UpdateAlbumParticipantAssignedCardFn: func(ctx context.Context, arg sqlc.UpdateAlbumParticipantAssignedCardParams) (sqlc.AlbumParticipant, error) {
			return sqlc.AlbumParticipant{}, errors.New("db error")
		},
	}
	_, err := admin_application.AssignCardToUser(mock, "u1", "c9")
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}
