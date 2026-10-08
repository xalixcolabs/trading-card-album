package admin_application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"com.xalixcolabs.trading-card-album/context/admin/model/dto"
	card_application "com.xalixcolabs.trading-card-album/context/card/application"
	card_dto "com.xalixcolabs.trading-card-album/context/card/model/dto"
	card_model "com.xalixcolabs.trading-card-album/context/card/model"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func ListCards(q database.Querier) ([]card_model.Card, error) {
	ctx := context.Background()
	cards, err := q.ListCards(ctx)
	if err != nil {
		return nil, err
	}
	return card_model.NewCardSliceFromSqlcCardSlice(cards), nil
}

func ListCardsByAlbum(q database.Querier, albumId string) ([]card_model.Card, error) {
	ctx := context.Background()
	cards, err := q.ListCardsByAlbumId(ctx, albumId)
	if err != nil {
		return nil, err
	}
	return card_model.NewCardSliceFromSqlcCardSlice(cards), nil
}

func DeleteCard(q database.Querier, cardId string) error {
	ctx := context.Background()
	return q.DeleteCard(ctx, cardId)
}

func CreateCard(q database.Querier, request admin_dto.CreateCardRequest) (card_model.Card, error) {
	ctx := context.Background()
	album, err := q.GetAlbum(ctx, request.AlbumId)
	if errors.Is(err, sql.ErrNoRows) {
		return card_model.Card{}, fmt.Errorf("Album no encontrado")
	}
	if err != nil {
		return card_model.Card{}, err
	}
	card, err := card_application.CreateCard(q, card_dto.CreateCardRequest{
		AlbumId:        album.ID,
		Number:         request.Number,
		Name:           request.Name,
		Description:    request.Description,
		ImageUrl:       request.ImageUrl,
		AutoAssignable: request.AutoAssignable,
	})
	if err != nil {
		return card_model.Card{}, err
	}
	if _, err = q.CreateCardPoolRow(ctx, sqlc.CreateCardPoolRowParams{
		AlbumID: album.ID,
		CardID:  card.ID,
	}); err != nil {
		return card_model.Card{}, err
	}
	if _, err = q.UpdateAlbum(ctx, sqlc.UpdateAlbumParams{
		Title:      album.Title,
		TotalCards: album.TotalCards + 1,
		ID:         album.ID,
	}); err != nil {
		return card_model.Card{}, err
	}
	return card, nil
}

func UpdateCard(q database.Querier, cardId string, request admin_dto.UpdateCardRequest) (card_model.Card, error) {
	ctx := context.Background()
	card, err := q.GetCard(ctx, cardId)
	if errors.Is(err, sql.ErrNoRows) {
		return card_model.Card{}, err
	}
	if err != nil {
		return card_model.Card{}, err
	}
	// Si no se especifica, se conserva la auto asignabilidad actual.
	autoAssignable := card.AutoAssignable != 0
	if request.AutoAssignable != nil {
		autoAssignable = *request.AutoAssignable
	}
	autoAssignableInt := int64(0)
	if autoAssignable {
		autoAssignableInt = 1
	}
	updated, err := q.UpdateCard(ctx, sqlc.UpdateCardParams{
		AlbumID:        card.AlbumID,
		Number:         request.Number,
		Name:           request.Name,
		Description:    request.Description,
		ImageUrl:       request.ImageUrl,
		AutoAssignable: autoAssignableInt,
		UpdatedAt:      time.Now().Unix(),
		ID:             cardId,
	})
	if err != nil {
		return card_model.Card{}, err
	}
	return card_model.NewCardFromSqlcCard(updated), nil
}