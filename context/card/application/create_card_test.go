package card_application_test

import (
	"context"
	"errors"
	"testing"

	"com.xalixcolabs.trading-card-album/context/card/application"
	dto "com.xalixcolabs.trading-card-album/context/card/model/dto"
	card_model "com.xalixcolabs.trading-card-album/context/card/model"
	"com.xalixcolabs.trading-card-album/database/queriermock"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func TestCreateCardReturnsCardModel(t *testing.T) {
	mock := &queriermock.Querier{
		CreateCardFn: func(ctx context.Context, arg sqlc.CreateCardParams) (sqlc.Card, error) {
			return sqlc.Card{
				ID:          arg.ID,
				AlbumID:     arg.AlbumID,
				Number:      arg.Number,
				Name:        arg.Name,
				Description: arg.Description,
				ImageUrl:    arg.ImageUrl,
				CreatedAt:   arg.CreatedAt,
				UpdatedAt:   arg.UpdatedAt,
			}, nil
		},
	}

	request := dto.CreateCardRequest{
		AlbumId:     "album-1",
		Number:      "01",
		Name:        "Gopher",
		Description: "A Go card",
		ImageUrl:    "http://localhost:8081/gopher.webp",
	}

	card, err := card_application.CreateCard(mock, request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if card.ID == "" {
		t.Error("expected a generated card id")
	}
	if card.AlbumId != "album-1" || card.Number != "01" || card.Name != "Gopher" {
		t.Errorf("unexpected card: %+v", card)
	}
	expected := card_model.NewCardFromSqlcCard(sqlc.Card{
		ID:          card.ID,
		AlbumID:     "album-1",
		Number:      "01",
		Name:        "Gopher",
		Description: "A Go card",
		ImageUrl:    "http://localhost:8081/gopher.webp",
		CreatedAt:   card.CreatedAt,
		UpdatedAt:   card.UpdatedAt,
	})
	if card != expected {
		t.Errorf("expected card %+v, got %+v", expected, card)
	}
}

func TestCreateCardReturnsError(t *testing.T) {
	mock := &queriermock.Querier{
		CreateCardFn: func(ctx context.Context, arg sqlc.CreateCardParams) (sqlc.Card, error) {
			return sqlc.Card{}, errors.New("db error")
		},
	}
	_, err := card_application.CreateCard(mock, dto.CreateCardRequest{AlbumId: "album-1"})
	if err == nil || err.Error() != "db error" {
		t.Errorf("expected db error, got %v", err)
	}
}