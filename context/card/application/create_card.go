package card_application

import (
	"context"
	"time"

	"com.xalixcolabs.trading-card-album/context/card/model"
	"com.xalixcolabs.trading-card-album/context/card/model/dto"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
	"github.com/matoous/go-nanoid/v2"
)

func CreateCard(q database.Querier, request card_dto.CreateCardRequest) (card_model.Card, error) {
	ctx := context.Background()
	id, _ := gonanoid.New()
	// Por defecto las tarjetas son auto asignables; solo se desactiva si se pide.
	autoAssignable := int64(1)
	if request.AutoAssignable != nil && !*request.AutoAssignable {
		autoAssignable = 0
	}
	card, err := q.CreateCard(ctx, sqlc.CreateCardParams{
		ID:             id,
		AlbumID:        request.AlbumId,
		Number:         request.Number,
		Name:           request.Name,
		Description:    request.Description,
		ImageUrl:       request.ImageUrl,
		AutoAssignable: autoAssignable,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
	})
	return card_model.NewCardFromSqlcCard(card), err
}
