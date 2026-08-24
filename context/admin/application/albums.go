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
	return q.DeleteAlbum(ctx, albumId)
}