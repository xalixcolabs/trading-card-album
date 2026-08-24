package admin_application

import (
	"context"
	"database/sql"
	"errors"

	"com.xalixcolabs.trading-card-album/context/admin/model/dto"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func UpdateAlbum(q database.Querier, albumId string, request admin_dto.UpdateAlbumRequest) (admin_dto.Album, error) {
	ctx := context.Background()
	album, err := q.GetAlbum(ctx, albumId)
	if errors.Is(err, sql.ErrNoRows) {
		return admin_dto.Album{}, err
	}
	if err != nil {
		return admin_dto.Album{}, err
	}
	updated, err := q.UpdateAlbum(ctx, sqlc.UpdateAlbumParams{
		Title:      request.Title,
		TotalCards: album.TotalCards,
		ID:         albumId,
	})
	if err != nil {
		return admin_dto.Album{}, err
	}
	return admin_dto.Album{
		ID:         updated.ID,
		Title:      updated.Title,
		TotalCards: updated.TotalCards,
		CreatedAt:  updated.CreatedAt,
	}, nil
}