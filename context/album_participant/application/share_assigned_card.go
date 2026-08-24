package album_participant_application

import (
	"context"

	"com.xalixcolabs.trading-card-album/context/album_participant/model/dto"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func ShareAssignedCard(q database.Querier, user user_model.User, albumId string) (album_participant_dto.ShareAssignedCardResponse, error) {
	ctx := context.Background()
	albumParticipant, err := q.GetAlbumParticipant(ctx, sqlc.GetAlbumParticipantParams{
		UserID:  user.ID,
		AlbumID: albumId,
	})
	if err != nil {
		return album_participant_dto.ShareAssignedCardResponse{}, err
	}
	return album_participant_dto.ShareAssignedCardResponse{
		ContactId: user.ID,
		AlbumId:   albumId,
		CardId:    albumParticipant.AssignedCardID,
		Secret:    albumParticipant.Secret,
	}, err
}
