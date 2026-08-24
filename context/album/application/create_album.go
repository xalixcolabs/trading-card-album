package album_application

import (
	"context"
	"time"

	"com.xalixcolabs.trading-card-album/context/album/model"
	"com.xalixcolabs.trading-card-album/context/album/model/dto"
	"com.xalixcolabs.trading-card-album/context/card/application"
	"com.xalixcolabs.trading-card-album/context/card/model"
	"com.xalixcolabs.trading-card-album/context/card/model/dto"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
	"github.com/matoous/go-nanoid/v2"
)

func CreateAlbum(q database.Querier, request album_dto.CreateAlbumRequest) (album_model.Album, error) {
	ctx := context.Background()
	id, _ := gonanoid.New()
	album, err := q.CreateAlbum(ctx, sqlc.CreateAlbumParams{
		ID:         id,
		Title:      request.Title,
		TotalCards: int64(len(request.Cards)),
		CreatedAt:  time.Now().Unix(),
	})
	var cards []card_model.Card
	for _, card := range request.Cards {
		newCard, err := card_application.CreateCard(q, card_dto.CreateCardRequest{
			AlbumId:     album.ID,
			Name:        card.Name,
			Description: card.Description,
			Number:      card.Number,
			ImageUrl:    card.ImageUrl,
		})
		if err != nil {
			return album_model.Album{}, err
		}
		_, err = q.CreateCardPoolRow(ctx, sqlc.CreateCardPoolRowParams{
			AlbumID: album.ID,
			CardID:  newCard.ID,
		})
		if err != nil {
			return album_model.Album{}, err
		}
		cards = append(cards, card_model.Card{
			ID:          newCard.ID,
			AlbumId:     newCard.AlbumId,
			Number:      newCard.Number,
			Name:        newCard.Name,
			Description: newCard.Description,
			ImageUrl:    newCard.ImageUrl,
			CreatedAt:   newCard.CreatedAt,
			UpdatedAt:   newCard.UpdatedAt,
		})
	}
	return album_model.NewAlbumFromSqlcAlbum(album, cards), err
}
