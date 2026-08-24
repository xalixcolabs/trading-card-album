package admin_application

import (
	"context"

	"com.xalixcolabs.trading-card-album/context/admin/model/dto"
	"com.xalixcolabs.trading-card-album/database"
)

func ListAlbums(q database.Querier) ([]admin_dto.Album, error) {
	ctx := context.Background()
	rows, err := q.ListAlbumsWithStats(ctx)
	if err != nil {
		return nil, err
	}
	albums := make([]admin_dto.Album, 0, len(rows))
	for _, row := range rows {
		albums = append(albums, admin_dto.Album{
			ID:               row.ID,
			Title:            row.Title,
			TotalCards:       row.TotalCards,
			CreatedAt:        row.CreatedAt,
			ParticipantCount: row.ParticipantCount,
		})
	}
	return albums, nil
}

func DeleteAlbum(q database.Querier, albumId string) error {
	ctx := context.Background()
	// SQLite no activa FK por defecto: se limpian los registros relacionados
	// de forma explícita para no dejar datos huérfanos.
	if err := q.DeleteUserCollectionByAlbumId(ctx, albumId); err != nil {
		return err
	}
	if err := q.DeleteCardPoolByAlbumId(ctx, albumId); err != nil {
		return err
	}
	if err := q.DeleteAlbumParticipantsByAlbumId(ctx, albumId); err != nil {
		return err
	}
	if err := q.DeleteCardsByAlbumId(ctx, albumId); err != nil {
		return err
	}
	return q.DeleteAlbum(ctx, albumId)
}