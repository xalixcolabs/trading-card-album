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
)

// GiftCardToUser regala una tarjeta de un álbum a un usuario: la suma a su
// colección sin cambiar su tarjeta asignada ni registrar un contacto.
func GiftCardToUser(q database.Querier, userId string, cardId string) (card_model.Card, error) {
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
