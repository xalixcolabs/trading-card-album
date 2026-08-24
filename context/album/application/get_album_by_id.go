package album_application

import (
	"context"
	"database/sql"
	"errors"

	"com.xalixcolabs.trading-card-album/context/album/model"
	card_model "com.xalixcolabs.trading-card-album/context/card/model"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func GetAlbumById(q database.Querier, user user_model.User, albumId string) (album_model.Album, error) {
	ctx := context.Background()
	album, err := q.GetAlbum(ctx, albumId)
	if errors.Is(err, sql.ErrNoRows) {
		return album_model.Album{}, err
	}
	if err != nil {
		return album_model.Album{}, err
	}
	// Solo los participantes del álbum pueden consultarlo.
	_, err = q.GetAlbumParticipant(ctx, sqlc.GetAlbumParticipantParams{
		AlbumID: album.ID,
		UserID:  user.ID,
	})
	if err != nil {
		return album_model.Album{}, err
	}
	// Solo se exponen las tarjetas recolectadas por el usuario; el resto
	// permanece oculto para no romper la sorpresa.
	collected, err := q.GetUserCollection(ctx, sqlc.GetUserCollectionParams{
		UserID:  user.ID,
		AlbumID: album.ID,
	})
	if err != nil {
		return album_model.Album{}, err
	}
	return album_model.NewAlbumFromSqlcAlbum(album, card_model.NewCardSliceFromSqlcCardSlice(collected)), nil
}