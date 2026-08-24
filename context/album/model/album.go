package album_model

import (
	card_model "com.xalixcolabs.trading-card-album/context/card/model"
	"com.xalixcolabs.trading-card-album/database/sqlc"
)

type Album struct {
	ID         string            `json:"id"`
	Title      string            `json:"title"`
	TotalCards int64             `json:"total_cards"`
	Cards      []card_model.Card `json:"cards"`
	CreatedAt  int64             `json:"created_at"`
}

func NewAlbumFromSqlcAlbum(album sqlc.Album, cards []card_model.Card) Album {
	return Album{
		ID:         album.ID,
		Title:      album.Title,
		TotalCards: album.TotalCards,
		Cards:      cards,
		CreatedAt:  album.CreatedAt,
	}
}
