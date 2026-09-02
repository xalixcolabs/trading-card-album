package album_application

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"com.xalixcolabs.trading-card-album/context/album/model/dto"
	"com.xalixcolabs.trading-card-album/context/card/model"
	"com.xalixcolabs.trading-card-album/context/events"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func RegisterCard(q database.Querier, user user_model.User, request album_dto.RegisterCardRequest) (card_model.Card, error) {
	ctx := context.Background()
	album, err := q.GetAlbum(ctx, request.AlbumId)
	if errors.Is(err, sql.ErrNoRows) {
		return card_model.Card{}, err
	}
	contact, err := q.GetUser(ctx, request.ContactId)
	if errors.Is(err, sql.ErrNoRows) {
		return card_model.Card{}, err
	}
	q.CreateContact(ctx, sqlc.CreateContactParams{
		UserID:    user.ID,
		MetUserID: request.ContactId,
		ScannedAt: time.Now().Unix(),
	})
	albumParticipant, err := q.GetAlbumParticipant(ctx, sqlc.GetAlbumParticipantParams{
		AlbumID: album.ID,
		UserID:  contact.ID,
	})
	if err != nil {
		return card_model.Card{}, err
	}
	if albumParticipant.Secret != request.Secret {
		return card_model.Card{}, errors.New("Esta QR ya ha sido escaneado")
	}
	owned, err := q.CardInCollection(ctx, sqlc.CardInCollectionParams{
		UserID:  user.ID,
		AlbumID: album.ID,
		CardID:  albumParticipant.AssignedCardID,
	})
	if err != nil {
		return card_model.Card{}, err
	}
	if !owned {
		if err = q.CollectCard(ctx, sqlc.CollectCardParams{
			UserID:     user.ID,
			AlbumID:    album.ID,
			CardID:     albumParticipant.AssignedCardID,
			UnlockedAt: time.Now().Unix(),
		}); err != nil {
			return card_model.Card{}, err
		}
	}
	secret, _ := gonanoid.New()
	_, err = q.UpdateAlbumParticipantSecret(ctx, sqlc.UpdateAlbumParticipantSecretParams{
		Secret:  secret,
		AlbumID: album.ID,
		UserID:  contact.ID,
	})
	if err != nil {
		return card_model.Card{}, err
	}
	// El QR del participante cambió (secret rotado): avisar por SSE para que
	// su cliente refresque el QR.
	eventData, _ := json.Marshal(map[string]string{"album_id": album.ID})
	events.Publish(contact.ID, eventData)
	card, _ := q.GetCard(ctx, albumParticipant.AssignedCardID)
	return card_model.NewCardFromSqlcCard(card), nil
}
