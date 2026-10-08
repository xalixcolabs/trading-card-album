package admin_application_test

import (
	"context"
	"testing"

	"com.xalixcolabs.trading-card-album/context/admin/application"
	dto "com.xalixcolabs.trading-card-album/context/admin/model/dto"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestAdminCreateCardPassesAutoAssignable(t *testing.T) {
	no := false
	var captured sqlc.CreateCardParams
	mock := &queriermock.Querier{
		GetAlbumFn: func(ctx context.Context, id string) (sqlc.Album, error) {
			return sqlc.Album{ID: "album-1", Title: "Album", TotalCards: 2}, nil
		},
		CreateCardFn: func(ctx context.Context, arg sqlc.CreateCardParams) (sqlc.Card, error) {
			captured = arg
			return sqlc.Card{ID: arg.ID, AlbumID: arg.AlbumID, AutoAssignable: arg.AutoAssignable}, nil
		},
		CreateCardPoolRowFn: func(ctx context.Context, arg sqlc.CreateCardPoolRowParams) (sqlc.CardPool, error) {
			return sqlc.CardPool{}, nil
		},
		UpdateAlbumFn: func(ctx context.Context, arg sqlc.UpdateAlbumParams) (sqlc.Album, error) {
			return sqlc.Album{}, nil
		},
	}
	if _, err := admin_application.CreateCard(mock, dto.CreateCardRequest{
		AlbumId:        "album-1",
		Number:         "01",
		Name:           "Gopher",
		ImageUrl:       "http://x/gopher.webp",
		AutoAssignable: &no,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.AutoAssignable != 0 {
		t.Errorf("expected auto_assignable 0, got %d", captured.AutoAssignable)
	}
}

func TestAdminUpdateCardPreservesAutoAssignableWhenOmitted(t *testing.T) {
	var captured sqlc.UpdateCardParams
	mock := &queriermock.Querier{
		GetCardFn: func(ctx context.Context, id string) (sqlc.Card, error) {
			return sqlc.Card{ID: "c1", AlbumID: "album-1", AutoAssignable: 0}, nil
		},
		UpdateCardFn: func(ctx context.Context, arg sqlc.UpdateCardParams) (sqlc.Card, error) {
			captured = arg
			return sqlc.Card{ID: arg.ID, AutoAssignable: arg.AutoAssignable}, nil
		},
	}
	if _, err := admin_application.UpdateCard(mock, "c1", dto.UpdateCardRequest{
		Number: "01", Name: "Gopher", ImageUrl: "http://x/gopher.webp",
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.AutoAssignable != 0 {
		t.Errorf("expected auto_assignable to stay 0, got %d", captured.AutoAssignable)
	}
}

func TestAdminUpdateCardTogglesAutoAssignable(t *testing.T) {
	yes := true
	var captured sqlc.UpdateCardParams
	mock := &queriermock.Querier{
		GetCardFn: func(ctx context.Context, id string) (sqlc.Card, error) {
			return sqlc.Card{ID: "c1", AlbumID: "album-1", AutoAssignable: 0}, nil
		},
		UpdateCardFn: func(ctx context.Context, arg sqlc.UpdateCardParams) (sqlc.Card, error) {
			captured = arg
			return sqlc.Card{ID: arg.ID, AutoAssignable: arg.AutoAssignable}, nil
		},
	}
	card, err := admin_application.UpdateCard(mock, "c1", dto.UpdateCardRequest{
		Number: "01", Name: "Gopher", ImageUrl: "http://x/gopher.webp", AutoAssignable: &yes,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.AutoAssignable != 1 {
		t.Errorf("expected auto_assignable 1, got %d", captured.AutoAssignable)
	}
	if !card.AutoAssignable {
		t.Error("expected returned card to be auto assignable")
	}
}
