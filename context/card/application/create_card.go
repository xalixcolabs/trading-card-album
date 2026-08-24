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

func CreateCard(request card_dto.CreateCardRequest) (card_model.Card, error) {
	ctx := context.Background()
	queries := sqlc.New(database.GetDatabase())
	id, _ := gonanoid.New()
	card, err := queries.CreateCard(ctx, sqlc.CreateCardParams{
		ID:          id,
		AlbumID:     request.AlbumId,
		Number:      request.Number,
		Name:        request.Name,
		Description: request.Description,
		ImageUrl:    request.ImageUrl,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	})
	return card_model.NewCardFromSqlcCard(card), err
}
