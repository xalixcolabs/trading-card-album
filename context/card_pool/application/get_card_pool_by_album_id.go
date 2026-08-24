package card_pool_application

import (
	"context"

	"com.xalixcolabs.trading-card-album/context/card/model"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func GetMyCardsByAlbumId(user user_model.User, albumId string) ([]card_model.Card, error) {
	ctx := context.Background()
	queries := sqlc.New(database.GetDatabase())
	cards, err := queries.GetUserCollection(ctx, sqlc.GetUserCollectionParams{
		UserID:  user.ID,
		AlbumID: albumId,
	})
	if err != nil {
		return []card_model.Card{}, nil
	}
	return card_model.NewCardSliceFromSqlcCardSlice(cards), nil
}
