package album_application_test

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"testing"

	"com.xalixcolabs.trading-card-album/context/album/application"
	"com.xalixcolabs.trading-card-album/context/album/model/dto"
	card_model "com.xalixcolabs.trading-card-album/context/card/model"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func testUser() user_model.User {
	return user_model.User{ID: "user-1", Email: "user@example.com"}
}

func testRegisterCardRequest() album_dto.RegisterCardRequest {
	return album_dto.RegisterCardRequest{
		AlbumId:   "album-1",
		ContactId: "contact-1",
		CardId:    "card-1",
		Secret:    "secret-abc",
	}
}

func testAlbumParticipant() sqlc.AlbumParticipant {
	return sqlc.AlbumParticipant{
		AlbumID:        "album-1",
		UserID:         "contact-1",
		AssignedCardID: "card-1",
		Secret:         "secret-abc",
	}
}

func testCard() sqlc.Card {
	return sqlc.Card{
		ID:          "card-1",
		AlbumID:     "album-1",
		Number:      "01",
		Name:        "Test Card",
		Description: "desc",
		ImageUrl:    "http://localhost:8081/card-1.webp",
	}
}

func TestRegisterCardRegistersCardAndContactWhenNotOwned(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1", Title: "Album"}, nil
		},
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "contact-1"}, nil
		},
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return testAlbumParticipant(), nil
		},
		CardInCollectionFn: func(ctx context.Context, arg sqlc.CardInCollectionParams) (bool, error) {
			return false, nil
		},
		GetCardFn: func(ctx context.Context, id string) (sqlc.Card, error) {
			return testCard(), nil
		},
	}

	contactCreated := false
	mock.CreateContactFn = func(ctx context.Context, arg sqlc.CreateContactParams) (sqlc.Contact, error) {
		contactCreated = true
		return sqlc.Contact{UserID: arg.UserID, MetUserID: arg.MetUserID}, nil
	}

	cardCollected := false
	mock.CollectCardFn = func(ctx context.Context, arg sqlc.CollectCardParams) error {
		cardCollected = true
		return nil
	}

	secretRotated := false
	mock.UpdateAlbumParticipantSecretFn = func(ctx context.Context, arg sqlc.UpdateAlbumParticipantSecretParams) (sqlc.AlbumParticipant, error) {
		secretRotated = true
		return testAlbumParticipant(), nil
	}

	card, err := album_application.RegisterCard(mock, testUser(), testRegisterCardRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contactCreated {
		t.Error("expected contact to be registered")
	}
	if !cardCollected {
		t.Error("expected card to be collected when not owned")
	}
	if !secretRotated {
		t.Error("expected participant secret to be rotated")
	}
	if card.ID != "card-1" {
		t.Errorf("expected card id card-1, got %s", card.ID)
	}
}

func TestRegisterCardRegistersOnlyContactWhenCardAlreadyOwned(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1", Title: "Album"}, nil
		},
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "contact-1"}, nil
		},
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return testAlbumParticipant(), nil
		},
		CardInCollectionFn: func(ctx context.Context, arg sqlc.CardInCollectionParams) (bool, error) {
			return true, nil
		},
		GetCardFn: func(ctx context.Context, id string) (sqlc.Card, error) {
			return testCard(), nil
		},
	}

	contactCreated := false
	mock.CreateContactFn = func(ctx context.Context, arg sqlc.CreateContactParams) (sqlc.Contact, error) {
		contactCreated = true
		return sqlc.Contact{}, nil
	}

	cardCollected := false
	mock.CollectCardFn = func(ctx context.Context, arg sqlc.CollectCardParams) error {
		cardCollected = true
		return nil
	}

	card, err := album_application.RegisterCard(mock, testUser(), testRegisterCardRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contactCreated {
		t.Error("expected contact to be registered even when card already owned")
	}
	if cardCollected {
		t.Error("expected card NOT to be collected when already owned")
	}
	if card.ID != "card-1" {
		t.Errorf("expected card id card-1, got %s", card.ID)
	}
}

func TestRegisterCardAutoJoinsAlbumWhenNotParticipant(t *testing.T) {
	var joinedUser string
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1", Title: "Album", TotalCards: 3}, nil
		},
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "contact-1"}, nil
		},
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			if arg.UserID == "contact-1" {
				return testAlbumParticipant(), nil
			}
			return sqlc.AlbumParticipant{}, sql.ErrNoRows
		},
		GetRandomAvailableCardFn: func(ctx context.Context, albumID string) (string, error) {
			return "card-9", nil
		},
		MarkCardAsDrawnFn: func(ctx context.Context, arg sqlc.MarkCardAsDrawnParams) (sqlc.CardPool, error) {
			return sqlc.CardPool{AlbumID: arg.AlbumID, CardID: arg.CardID, IsDrawn: 1}, nil
		},
		CreateAlbumParticipantFn: func(ctx context.Context, arg sqlc.CreateAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			joinedUser = arg.UserID
			return sqlc.AlbumParticipant{
				AlbumID:        arg.AlbumID,
				UserID:         arg.UserID,
				AssignedCardID: arg.AssignedCardID,
				JoinedAt:       arg.JoinedAt,
				Secret:         arg.Secret,
			}, nil
		},
		CardInCollectionFn: func(ctx context.Context, arg sqlc.CardInCollectionParams) (bool, error) {
			return false, nil
		},
		GetCardFn: func(ctx context.Context, id string) (sqlc.Card, error) {
			return testCard(), nil
		},
	}

	mock.CreateContactFn = func(ctx context.Context, arg sqlc.CreateContactParams) (sqlc.Contact, error) {
		return sqlc.Contact{UserID: arg.UserID, MetUserID: arg.MetUserID}, nil
	}
	mock.UpdateAlbumParticipantSecretFn = func(ctx context.Context, arg sqlc.UpdateAlbumParticipantSecretParams) (sqlc.AlbumParticipant, error) {
		return testAlbumParticipant(), nil
	}

	collected := []string{}
	mock.CollectCardFn = func(ctx context.Context, arg sqlc.CollectCardParams) error {
		collected = append(collected, arg.CardID)
		return nil
	}

	card, err := album_application.RegisterCard(mock, testUser(), testRegisterCardRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if joinedUser != "user-1" {
		t.Errorf("expected scanner to auto-join the album, got %q", joinedUser)
	}
	if !slices.Contains(collected, "card-9") {
		t.Errorf("expected the scanner's own assigned card to be collected, got %v", collected)
	}
	if !slices.Contains(collected, "card-1") {
		t.Errorf("expected the scanned card to be collected, got %v", collected)
	}
	if card.ID != "card-1" {
		t.Errorf("expected card id card-1, got %s", card.ID)
	}
}

func TestRegisterCardDoesNotAssignCardWhenAlreadyParticipant(t *testing.T) {
	cardAssigned := false
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1", Title: "Album", TotalCards: 3}, nil
		},
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "contact-1"}, nil
		},
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return testAlbumParticipant(), nil
		},
		GetRandomAvailableCardFn: func(ctx context.Context, albumID string) (string, error) {
			cardAssigned = true
			return "card-9", nil
		},
		CardInCollectionFn: func(ctx context.Context, arg sqlc.CardInCollectionParams) (bool, error) {
			return true, nil
		},
		GetCardFn: func(ctx context.Context, id string) (sqlc.Card, error) {
			return testCard(), nil
		},
	}

	mock.CreateContactFn = func(ctx context.Context, arg sqlc.CreateContactParams) (sqlc.Contact, error) {
		return sqlc.Contact{}, nil
	}
	mock.UpdateAlbumParticipantSecretFn = func(ctx context.Context, arg sqlc.UpdateAlbumParticipantSecretParams) (sqlc.AlbumParticipant, error) {
		return testAlbumParticipant(), nil
	}

	_, err := album_application.RegisterCard(mock, testUser(), testRegisterCardRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cardAssigned {
		t.Error("expected NO card to be assigned when the scanner already participates")
	}
}

func TestRegisterCardReturnsErrorWhenAlbumNotFound(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{}, sql.ErrNoRows
		},
	}
	_, err := album_application.RegisterCard(mock, testUser(), testRegisterCardRequest())
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestRegisterCardReturnsErrorWhenContactNotFound(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1"}, nil
		},
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{}, sql.ErrNoRows
		},
	}
	_, err := album_application.RegisterCard(mock, testUser(), testRegisterCardRequest())
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestRegisterCardReturnsErrorOnSecretMismatch(t *testing.T) {
	participant := testAlbumParticipant()
	participant.Secret = "different-secret"

	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1"}, nil
		},
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "contact-1"}, nil
		},
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return participant, nil
		},
	}
	_, err := album_application.RegisterCard(mock, testUser(), testRegisterCardRequest())
	if err == nil || err.Error() != "Esta QR ya ha sido escaneado" {
		t.Errorf("expected 'Esta QR ya ha sido escaneado' error, got %v", err)
	}
}

func TestRegisterCardReturnsErrorWhenParticipantNotFound(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1"}, nil
		},
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "contact-1"}, nil
		},
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return sqlc.AlbumParticipant{}, sql.ErrNoRows
		},
	}
	_, err := album_application.RegisterCard(mock, testUser(), testRegisterCardRequest())
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestRegisterCardReturnsErrorWhenCardInCollectionFails(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1"}, nil
		},
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "contact-1"}, nil
		},
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return testAlbumParticipant(), nil
		},
		CardInCollectionFn: func(ctx context.Context, arg sqlc.CardInCollectionParams) (bool, error) {
			return false, errors.New("db error")
		},
	}
	_, err := album_application.RegisterCard(mock, testUser(), testRegisterCardRequest())
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestRegisterCardReturnsCardModel(t *testing.T) {
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1"}, nil
		},
		GetUserFn: func(ctx context.Context, id string) (sqlc.User, error) {
			return sqlc.User{ID: "contact-1"}, nil
		},
		GetAlbumParticipantFn: func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
			return testAlbumParticipant(), nil
		},
		CardInCollectionFn: func(ctx context.Context, arg sqlc.CardInCollectionParams) (bool, error) {
			return true, nil
		},
		GetCardFn: func(ctx context.Context, id string) (sqlc.Card, error) {
			return testCard(), nil
		},
	}

	card, err := album_application.RegisterCard(mock, testUser(), testRegisterCardRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := card_model.NewCardFromSqlcCard(testCard())
	if card != expected {
		t.Errorf("expected card %+v, got %+v", expected, card)
	}
}