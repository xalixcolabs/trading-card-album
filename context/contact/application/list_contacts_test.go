package contact_application_test

import (
	"context"
	"errors"
	"testing"

	"com.xalixcolabs.trading-card-album/context/contact/application"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestListContactsHidesPrivateEmail(t *testing.T) {
	mock := &queriermock.Querier{
		ListContactsFn: func(ctx context.Context, userID string) ([]sqlc.ListContactsRow, error) {
			return []sqlc.ListContactsRow{
				{ID: "u1", Name: "Privado", Email: "privado@example.com", PublicEmail: 0, PublicContact: "@privado"},
				{ID: "u2", Name: "Público", Email: "publico@example.com", PublicEmail: 1},
			}, nil
		},
	}

	contacts, err := contact_application.ListContacts(mock, user_model.User{ID: "me"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(contacts) != 2 {
		t.Fatalf("expected 2 contacts, got %d", len(contacts))
	}
	if contacts[0].Email != "" || contacts[0].Contact != "@privado" {
		t.Errorf("private contact should hide email and expose alt contact: %+v", contacts[0])
	}
	if contacts[1].Email != "publico@example.com" || contacts[1].Contact != "" {
		t.Errorf("public contact should expose email: %+v", contacts[1])
	}
}

func TestListContactsReturnsError(t *testing.T) {
	mock := &queriermock.Querier{
		ListContactsFn: func(ctx context.Context, userID string) ([]sqlc.ListContactsRow, error) {
			return nil, errors.New("db error")
		},
	}
	_, err := contact_application.ListContacts(mock, user_model.User{ID: "me"})
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}
