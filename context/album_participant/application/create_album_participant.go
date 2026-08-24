package album_participant_application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"com.xalixcolabs.trading-card-album/context/album_participant/model"
	"com.xalixcolabs.trading-card-album/context/album_participant/model/dto"
	"com.xalixcolabs.trading-card-album/context/card_pool/application"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
	"github.com/matoous/go-nanoid/v2"
)

func CreateAlbumParticipant(user user_model.User, request album_participant_dto.CreateAlbumParticipantRequest) (album_participant_model.AlbumParticipant, error) {
	ctx := context.Background()
	queries := sqlc.New(database.GetDatabase())
	album, err := queries.GetAlbum(ctx, request.AlbumId)
	if errors.Is(err, sql.ErrNoRows) {
		return album_participant_model.AlbumParticipant{}, fmt.Errorf("Album no encontrado")
	}
	secret, _ := gonanoid.New()
	cardId, err := card_pool_application.AssignCard(album.ID)
	if err != nil {
		return album_participant_model.AlbumParticipant{}, err
	}
	albumParticipant, err := queries.CreateAlbumParticipant(ctx, sqlc.CreateAlbumParticipantParams{
		AlbumID:        album.ID,
		UserID:         user.ID,
		AssignedCardID: cardId,
		Secret:         secret,
		JoinedAt:       time.Now().Unix(),
	})
	if err = queries.CollectCard(ctx, sqlc.CollectCardParams{
		UserID:     user.ID,
		AlbumID:    album.ID,
		CardID:     cardId,
		UnlockedAt: time.Now().Unix(),
	}); err != nil {
		return album_participant_model.AlbumParticipant{}, err
	}
	return album_participant_model.NewAlbumParticipantFromSqlcAlbumParticipant(albumParticipant),
		err
}
