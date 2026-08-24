package card_pool_application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"com.xalixcolabs.trading-card-album/database"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

func AssignCard(albumId string) (string, error) {
	ctx := context.Background()
	queries := sqlc.New(database.GetDatabase())
	cardId, err := queries.GetRandomAvailableCard(ctx, albumId)

	if errors.Is(err, sql.ErrNoRows) {
		if _, err := queries.ResetCardPool(ctx, albumId); err != nil {
			return "", fmt.Errorf("error al reiniciar el pool: %w", err)
		}

		cardId, err = queries.GetRandomAvailableCard(ctx, albumId)
		if err != nil {
			return "", fmt.Errorf("pool vacío tras reiniciar: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("error al buscar tarjeta: %w", err)
	}
	_, err = queries.MarkCardAsDrawn(ctx, sqlc.MarkCardAsDrawnParams{
		AlbumID: albumId,
		CardID:  cardId,
	})
	if err != nil {
		return "", fmt.Errorf("error al marcar tarjeta como asignada: %w", err)
	}
	return cardId, nil
}
