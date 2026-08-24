package album_participant_application_test

import (
	"context"
	"database/sql"
	"testing"

	"com.xalixcolabs.trading-card-album/context/album_participant/application"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestShareAssignedCardReturnsSharePayload(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return sqlc.AlbumParticipant{
				AlbumID:        "album-1",
				UserID:         "user-1",
				AssignedCardID: "card-7",
				Secret:         "top-secret",
			}, nil
		},
	}

	user := user_model.User{ID: "user-1"}
	response, err := album_participant_application.ShareAssignedCard(mock, user, "album-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.ContactId != "user-1" {
		t.Errorf("expected contact id user-1, got %s", response.ContactId)
	}
	if response.AlbumId != "album-1" {
		t.Errorf("expected album id album-1, got %s", response.AlbumId)
	}
	if response.CardId != "card-7" {
		t.Errorf("expected card id card-7, got %s", response.CardId)
	}
	if response.Secret != "top-secret" {
		t.Errorf("expected secret top-secret, got %s", response.Secret)
	}
}

func TestShareAssignedCardReturnsErrorWhenNotParticipant(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return sqlc.AlbumParticipant{}, sql.ErrNoRows
		},
	}
	_, err := album_participant_application.ShareAssignedCard(mock, user_model.User{ID: "user-1"}, "album-1")
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}