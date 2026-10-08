package admin_application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	card_model "com.xalixcolabs.trading-card-album/context/card/model"
	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
	"github.com/matoous/go-nanoid/v2"
)

// AssignCardToUser asigna una tarjeta específica a un usuario. Si el usuario no
// participa todavía en el álbum de la tarjeta, lo inscribe con esa tarjeta como
// su tarjeta asignada; si ya participa, reemplaza su tarjeta asignada. En ambos
// casos la tarjeta se marca como repartida y se suma a su colección.
func AssignCardToUser(q database.Querier, userId string, cardId string) (card_model.Card, error) {
	ctx := context.Background()
	if _, err := q.GetUser(ctx, userId); errors.Is(err, sql.ErrNoRows) {
		return card_model.Card{}, fmt.Errorf("Usuario no encontrado")
	} else if err != nil {
		return card_model.Card{}, err
	}

	card, err := q.GetCard(ctx, cardId)
	if errors.Is(err, sql.ErrNoRows) {
		return card_model.Card{}, fmt.Errorf("Tarjeta no encontrada")
	}
	if err != nil {
		return card_model.Card{}, err
	}

	_, err = q.GetAlbumParticipant(ctx, sqlc.GetAlbumParticipantParams{
		AlbumID: card.AlbumID,
		UserID:  userId,
	})
	switch {
	case errors.Is(err, sql.ErrNoRows):
		secret, _ := gonanoid.New()
		if _, err = q.CreateAlbumParticipant(ctx, sqlc.CreateAlbumParticipantParams{
			AlbumID:        card.AlbumID,
			UserID:         userId,
			AssignedCardID: cardId,
			Secret:         secret,
			JoinedAt:       time.Now().Unix(),
		}); err != nil {
			return card_model.Card{}, err
		}
	case err != nil:
		return card_model.Card{}, err
	default:
		if _, err = q.UpdateAlbumParticipantAssignedCard(ctx, sqlc.UpdateAlbumParticipantAssignedCardParams{
			AssignedCardID: cardId,
			AlbumID:        card.AlbumID,
			UserID:         userId,
		}); err != nil {
			return card_model.Card{}, err
		}
	}

	if _, err = q.MarkCardAsDrawn(ctx, sqlc.MarkCardAsDrawnParams{
		AlbumID: card.AlbumID,
		CardID:  cardId,
	}); err != nil {
		return card_model.Card{}, err
	}
	if err = q.CollectCard(ctx, sqlc.CollectCardParams{
		UserID:     userId,
		AlbumID:    card.AlbumID,
		CardID:     cardId,
		UnlockedAt: time.Now().Unix(),
	}); err != nil {
		return card_model.Card{}, err
	}

	return card_model.NewCardFromSqlcCard(card), nil
}
