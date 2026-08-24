package card_application

import (
	"context"

	"com.xalixcolabs.trading-card-album/context/card/model"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func GetCardByAlbumId(user user_model.User, albumId string) (card_model.Card, error) {
	ctx := context.Background()
	queries := sqlc.New(database.GetDatabase())
	albumParticipant, err := queries.GetAlbumParticipant(ctx, sqlc.GetAlbumParticipantParams{
		UserID:  user.ID,
		AlbumID: albumId,
	})
	if err != nil {
		return card_model.Card{}, nil
	}
	card, err := queries.GetCardByAlbumParticipant(ctx, sqlc.GetCardByAlbumParticipantParams{
		UserID:         user.ID,
		AlbumID:        albumId,
		AssignedCardID: albumParticipant.AssignedCardID,
	})
	if err != nil {
		return card_model.Card{}, nil
	}
	return card_model.NewCardFromSqlcCard(card), nil
}
