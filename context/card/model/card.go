package card_model

import "com.xalixcolabs.trading-card-album/database/sqlc"

type Card struct {
	ID          string `json:"id"`
	AlbumId     string `json:"album_id"`
	Number      string `json:"number"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

func NewCardFromSqlcCard(card sqlc.Card) Card {
	return Card{
		ID:          card.ID,
		AlbumId:     card.AlbumID,
		Number:      card.Number,
		Name:        card.Name,
		Description: card.Description,
		ImageUrl:    card.ImageUrl,
		CreatedAt:   card.CreatedAt,
		UpdatedAt:   card.UpdatedAt,
	}
}

func NewCardSliceFromSqlcCardSlice(cards []sqlc.Card) []Card {
	cards_ := make([]Card, 0)
	for _, card := range cards {
		cards_ = append(cards_, Card{
			ID:          card.ID,
			AlbumId:     card.AlbumID,
			Number:      card.Number,
			Name:        card.Name,
			Description: card.Description,
			ImageUrl:    card.ImageUrl,
			CreatedAt:   card.CreatedAt,
			UpdatedAt:   card.UpdatedAt,
		})
	}
	return cards_
}
