package admin_application

import (
	"context"

	card_model "com.xalixcolabs.trading-card-album/context/card/model"
	"com.xalixcolabs.trading-card-album/database"
)

func ListCards(q database.Querier) ([]card_model.Card, error) {
	ctx := context.Background()
	cards, err := q.ListCards(ctx)
	if err != nil {
		return nil, err
	}
	return card_model.NewCardSliceFromSqlcCardSlice(cards), nil
}

func DeleteCard(q database.Querier, cardId string) error {
	ctx := context.Background()
	return q.DeleteCard(ctx, cardId)
}