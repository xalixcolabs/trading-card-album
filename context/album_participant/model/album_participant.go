package album_participant_model

import "com.xalixcolabs.trading-card-album/database/sqlc"

type AlbumParticipant struct {
	AlbumID        string `json:"album_id"`
	UserID         string `json:"user_id"`
	AssignedCardID string `json:"assigned_card_id"`
	JoinedAt       int64  `json:"joined_at"`
	Secret         string `json:"-"`
}

func NewAlbumParticipantFromSqlcAlbumParticipant(albumParticipant sqlc.AlbumParticipant) AlbumParticipant {
	return AlbumParticipant{
		AlbumID:        albumParticipant.AlbumID,
		UserID:         albumParticipant.UserID,
		AssignedCardID: albumParticipant.AssignedCardID,
		JoinedAt:       albumParticipant.JoinedAt,
		Secret:         albumParticipant.Secret,
	}
}
