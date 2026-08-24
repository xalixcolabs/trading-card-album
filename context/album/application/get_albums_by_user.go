package album_application

import (
	"context"

	"com.xalixcolabs.trading-card-album/context/album/model"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func GetAlbumsByUser(user user_model.User) ([]album_model.Album, error) {
	ctx := context.Background()
	queries := sqlc.New(database.GetDatabase())
	albums, err := queries.ListAlbumsByUserId(ctx, user.ID)
	albumsResponse := make([]album_model.Album, 0)
	for _, album := range albums {
		albumsResponse = append(albumsResponse, album_model.Album{
			ID:         album.ID,
			Title:      album.Title,
			TotalCards: album.TotalCards,
			CreatedAt:  album.CreatedAt,
		})
	}
	return albumsResponse, err
}
