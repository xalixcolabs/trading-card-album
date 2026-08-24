package database

import (
	"context"

	"com.xalixcolabs.trading-card-album/database/sqlc"
)

// Querier es la interfaz de acceso a datos que usa la capa de application.
// La implementa *sqlc.Queries y se inyecta para poder testear los casos de
// uso con un mock en lugar de tocar la base de datos real.
type Querier interface {
	GetAlbum(ctx context.Context, id string) (sqlc.Album, error)
	CreateAlbum(ctx context.Context, arg sqlc.CreateAlbumParams) (sqlc.Album, error)
	ListAlbumsByUserId(ctx context.Context, userID string) ([]sqlc.Album, error)

	GetCard(ctx context.Context, id string) (sqlc.Card, error)
	CreateCard(ctx context.Context, arg sqlc.CreateCardParams) (sqlc.Card, error)
	GetCardByAlbumParticipant(ctx context.Context, arg sqlc.GetCardByAlbumParticipantParams) (sqlc.Card, error)

	GetUser(ctx context.Context, id string) (sqlc.User, error)
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)

	CreateContact(ctx context.Context, arg sqlc.CreateContactParams) (sqlc.Contact, error)

	CreateAlbumParticipant(ctx context.Context, arg sqlc.CreateAlbumParticipantParams) (sqlc.AlbumParticipant, error)
	GetAlbumParticipant(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error)
	UpdateAlbumParticipantSecret(ctx context.Context, arg sqlc.UpdateAlbumParticipantSecretParams) (sqlc.AlbumParticipant, error)

	CreateCardPoolRow(ctx context.Context, arg sqlc.CreateCardPoolRowParams) (sqlc.CardPool, error)
	GetRandomAvailableCard(ctx context.Context, albumID string) (string, error)
	ResetCardPool(ctx context.Context, albumID string) (sqlc.CardPool, error)
	MarkCardAsDrawn(ctx context.Context, arg sqlc.MarkCardAsDrawnParams) (sqlc.CardPool, error)

	CollectCard(ctx context.Context, arg sqlc.CollectCardParams) error
	GetUserCollection(ctx context.Context, arg sqlc.GetUserCollectionParams) ([]sqlc.Card, error)
	CardInCollection(ctx context.Context, arg sqlc.CardInCollectionParams) (bool, error)
}

// DefaultQuerier devuelve la implementación real de Querier usando la conexión
// global a la base de datos. Se usa en producción (resources).
func DefaultQuerier() Querier {
	return sqlc.New(GetDatabase())
}
