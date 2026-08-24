package album_participant_application_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"com.xalixcolabs.trading-card-album/context/album_participant/application"
	dto "com.xalixcolabs.trading-card-album/context/album_participant/model/dto"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestCreateAlbumParticipantAssignsCardAndCollects(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1", Title: "Album", TotalCards: 3}, nil
		},
		GetRandomAvailableCardFn: func(ctx context.Context, albumID string) (string, error) {
			return "card-9", nil
		},
		MarkCardAsDrawnFn: func(ctx context.Context, arg sqlc.MarkCardAsDrawnParams) (sqlc.CardPool, error) {
			return sqlc.CardPool{AlbumID: arg.AlbumID, CardID: arg.CardID, IsDrawn: 1}, nil
		},
		CreateAlbumParticipantFn: func(ctx context.Context, arg sqlc.CreateAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return sqlc.AlbumParticipant{
				AlbumID:        arg.AlbumID,
				UserID:         arg.UserID,
				AssignedCardID: arg.AssignedCardID,
				JoinedAt:       arg.JoinedAt,
				Secret:         arg.Secret,
			}, nil
		},
	}

	collected := false
	mock.CollectCardFn = func(ctx context.Context, arg sqlc.CollectCardParams) error {
		collected = true
		return nil
	}

	user := user_model.User{ID: "user-1"}
	request := dto.CreateAlbumParticipantRequest{AlbumId: "album-1"}

	participant, err := album_participant_application.CreateAlbumParticipant(mock, user, request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if participant.AlbumID != "album-1" {
		t.Errorf("expected album id album-1, got %s", participant.AlbumID)
	}
	if participant.UserID != "user-1" {
		t.Errorf("expected user id user-1, got %s", participant.UserID)
	}
	if participant.AssignedCardID != "card-9" {
		t.Errorf("expected assigned card card-9, got %s", participant.AssignedCardID)
	}
	if participant.Secret == "" {
		t.Error("expected a secret to be generated for the participant")
	}
	if !collected {
		t.Error("expected assigned card to be collected for the participant")
	}
}

func TestCreateAlbumParticipantReturnsErrorWhenAlbumNotFound(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{}, sql.ErrNoRows
		},
	}
	_, err := album_participant_application.CreateAlbumParticipant(
		mock, user_model.User{ID: "user-1"}, dto.CreateAlbumParticipantRequest{AlbumId: "album-1"},
	)
	if err == nil || err.Error() != "Album no encontrado" {
		t.Errorf("expected 'Album no encontrado' error, got %v", err)
	}
}

func TestCreateAlbumParticipantReturnsErrorWhenAssignCardFails(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1"}, nil
		},
		GetRandomAvailableCardFn: func(ctx context.Context, albumID string) (string, error) {
			return "", errors.New("pool error")
		},
	}
	_, err := album_participant_application.CreateAlbumParticipant(
		mock, user_model.User{ID: "user-1"}, dto.CreateAlbumParticipantRequest{AlbumId: "album-1"},
	)
	if err == nil || err.Error() != "error al buscar tarjeta: pool error" {
		t.Errorf("expected pool error, got %v", err)
	}
}