package queriermock

import (
	"context"

	"com.xalixcolabs.trading-card-album/database/sqlc"
)

// Querier es un mock de database.Querier para tests de la capa de application.
// Cada campo <X>Fn permite configurar el comportamiento de la llamada <X>; si
// no se configura, se devuelve el valor cero y nil como error.
type Querier struct {
	GetAlbumFn                     func(ctx context.Context, id string) (sqlc.Album, error)
	CreateAlbumFn                  func(ctx context.Context, arg sqlc.CreateAlbumParams) (sqlc.Album, error)
	ListAlbumsByUserIdFn           func(ctx context.Context, userID string) ([]sqlc.Album, error)
	ListAlbumsWithStatsFn          func(ctx context.Context) ([]sqlc.ListAlbumsWithStatsRow, error)
	DeleteAlbumFn                  func(ctx context.Context, id string) error
	GetOverviewStatsFn             func(ctx context.Context) (sqlc.GetOverviewStatsRow, error)
	GetCardFn                      func(ctx context.Context, id string) (sqlc.Card, error)
	CreateCardFn                   func(ctx context.Context, arg sqlc.CreateCardParams) (sqlc.Card, error)
	ListCardsFn                    func(ctx context.Context) ([]sqlc.Card, error)
	GetCardByAlbumParticipantFn    func(ctx context.Context, arg sqlc.GetCardByAlbumParticipantParams) (sqlc.Card, error)
	DeleteCardFn                   func(ctx context.Context, id string) error
	GetUserFn                      func(ctx context.Context, id string) (sqlc.User, error)
	CreateUserFn                   func(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	UpdateUserFn                   func(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)
	UpdateUserIsAdminFn            func(ctx context.Context, arg sqlc.UpdateUserIsAdminParams) (sqlc.User, error)
	ListUsersFn                    func(ctx context.Context) ([]sqlc.User, error)
CreateContactFn             func(ctx context.Context, arg sqlc.CreateContactParams) (sqlc.Contact, error)
	ListContactsFn              func(ctx context.Context, userID string) ([]sqlc.ListContactsRow, error)
	CreateAlbumParticipantFn       func(ctx context.Context, arg sqlc.CreateAlbumParticipantParams) (sqlc.AlbumParticipant, error)
	GetAlbumParticipantFn          func(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error)
	UpdateAlbumParticipantSecretFn func(ctx context.Context, arg sqlc.UpdateAlbumParticipantSecretParams) (sqlc.AlbumParticipant, error)
	CreateCardPoolRowFn            func(ctx context.Context, arg sqlc.CreateCardPoolRowParams) (sqlc.CardPool, error)
	GetRandomAvailableCardFn       func(ctx context.Context, albumID string) (string, error)
	ResetCardPoolFn                func(ctx context.Context, albumID string) (sqlc.CardPool, error)
	MarkCardAsDrawnFn              func(ctx context.Context, arg sqlc.MarkCardAsDrawnParams) (sqlc.CardPool, error)
	CollectCardFn                  func(ctx context.Context, arg sqlc.CollectCardParams) error
	GetUserCollectionFn            func(ctx context.Context, arg sqlc.GetUserCollectionParams) ([]sqlc.Card, error)
	CardInCollectionFn             func(ctx context.Context, arg sqlc.CardInCollectionParams) (bool, error)
}

func (m *Querier) GetAlbum(ctx context.Context, id string) (sqlc.Album, error) {
	if m.GetAlbumFn != nil {
		return m.GetAlbumFn(ctx, id)
	}
	return sqlc.Album{}, nil
}

func (m *Querier) ListAlbumsWithStats(ctx context.Context) ([]sqlc.ListAlbumsWithStatsRow, error) {
	if m.ListAlbumsWithStatsFn != nil {
		return m.ListAlbumsWithStatsFn(ctx)
	}
	return nil, nil
}

func (m *Querier) DeleteAlbum(ctx context.Context, id string) error {
	if m.DeleteAlbumFn != nil {
		return m.DeleteAlbumFn(ctx, id)
	}
	return nil
}

func (m *Querier) GetOverviewStats(ctx context.Context) (sqlc.GetOverviewStatsRow, error) {
	if m.GetOverviewStatsFn != nil {
		return m.GetOverviewStatsFn(ctx)
	}
	return sqlc.GetOverviewStatsRow{}, nil
}

func (m *Querier) CreateAlbum(ctx context.Context, arg sqlc.CreateAlbumParams) (sqlc.Album, error) {
	if m.CreateAlbumFn != nil {
		return m.CreateAlbumFn(ctx, arg)
	}
	return sqlc.Album{}, nil
}

func (m *Querier) ListAlbumsByUserId(ctx context.Context, userID string) ([]sqlc.Album, error) {
	if m.ListAlbumsByUserIdFn != nil {
		return m.ListAlbumsByUserIdFn(ctx, userID)
	}
	return nil, nil
}

func (m *Querier) GetCard(ctx context.Context, id string) (sqlc.Card, error) {
	if m.GetCardFn != nil {
		return m.GetCardFn(ctx, id)
	}
	return sqlc.Card{}, nil
}

func (m *Querier) ListCards(ctx context.Context) ([]sqlc.Card, error) {
	if m.ListCardsFn != nil {
		return m.ListCardsFn(ctx)
	}
	return nil, nil
}

func (m *Querier) DeleteCard(ctx context.Context, id string) error {
	if m.DeleteCardFn != nil {
		return m.DeleteCardFn(ctx, id)
	}
	return nil
}

func (m *Querier) CreateCard(ctx context.Context, arg sqlc.CreateCardParams) (sqlc.Card, error) {
	if m.CreateCardFn != nil {
		return m.CreateCardFn(ctx, arg)
	}
	return sqlc.Card{}, nil
}

func (m *Querier) GetCardByAlbumParticipant(ctx context.Context, arg sqlc.GetCardByAlbumParticipantParams) (sqlc.Card, error) {
	if m.GetCardByAlbumParticipantFn != nil {
		return m.GetCardByAlbumParticipantFn(ctx, arg)
	}
	return sqlc.Card{}, nil
}

func (m *Querier) GetUser(ctx context.Context, id string) (sqlc.User, error) {
	if m.GetUserFn != nil {
		return m.GetUserFn(ctx, id)
	}
	return sqlc.User{}, nil
}

func (m *Querier) CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, arg)
	}
	return sqlc.User{}, nil
}

func (m *Querier) UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
	if m.UpdateUserFn != nil {
		return m.UpdateUserFn(ctx, arg)
	}
	return sqlc.User{}, nil
}

func (m *Querier) UpdateUserIsAdmin(ctx context.Context, arg sqlc.UpdateUserIsAdminParams) (sqlc.User, error) {
	if m.UpdateUserIsAdminFn != nil {
		return m.UpdateUserIsAdminFn(ctx, arg)
	}
	return sqlc.User{}, nil
}

func (m *Querier) ListUsers(ctx context.Context) ([]sqlc.User, error) {
	if m.ListUsersFn != nil {
		return m.ListUsersFn(ctx)
	}
	return nil, nil
}

func (m *Querier) CreateContact(ctx context.Context, arg sqlc.CreateContactParams) (sqlc.Contact, error) {
	if m.CreateContactFn != nil {
		return m.CreateContactFn(ctx, arg)
	}
	return sqlc.Contact{}, nil
}

func (m *Querier) ListContacts(ctx context.Context, userID string) ([]sqlc.ListContactsRow, error) {
	if m.ListContactsFn != nil {
		return m.ListContactsFn(ctx, userID)
	}
	return nil, nil
}

func (m *Querier) CreateAlbumParticipant(ctx context.Context, arg sqlc.CreateAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
	if m.CreateAlbumParticipantFn != nil {
		return m.CreateAlbumParticipantFn(ctx, arg)
	}
	return sqlc.AlbumParticipant{}, nil
}

func (m *Querier) GetAlbumParticipant(ctx context.Context, arg sqlc.GetAlbumParticipantParams) (sqlc.AlbumParticipant, error) {
	if m.GetAlbumParticipantFn != nil {
		return m.GetAlbumParticipantFn(ctx, arg)
	}
	return sqlc.AlbumParticipant{}, nil
}

func (m *Querier) UpdateAlbumParticipantSecret(ctx context.Context, arg sqlc.UpdateAlbumParticipantSecretParams) (sqlc.AlbumParticipant, error) {
	if m.UpdateAlbumParticipantSecretFn != nil {
		return m.UpdateAlbumParticipantSecretFn(ctx, arg)
	}
	return sqlc.AlbumParticipant{}, nil
}

func (m *Querier) CreateCardPoolRow(ctx context.Context, arg sqlc.CreateCardPoolRowParams) (sqlc.CardPool, error) {
	if m.CreateCardPoolRowFn != nil {
		return m.CreateCardPoolRowFn(ctx, arg)
	}
	return sqlc.CardPool{}, nil
}

func (m *Querier) GetRandomAvailableCard(ctx context.Context, albumID string) (string, error) {
	if m.GetRandomAvailableCardFn != nil {
		return m.GetRandomAvailableCardFn(ctx, albumID)
	}
	return "", nil
}

func (m *Querier) ResetCardPool(ctx context.Context, albumID string) (sqlc.CardPool, error) {
	if m.ResetCardPoolFn != nil {
		return m.ResetCardPoolFn(ctx, albumID)
	}
	return sqlc.CardPool{}, nil
}

func (m *Querier) MarkCardAsDrawn(ctx context.Context, arg sqlc.MarkCardAsDrawnParams) (sqlc.CardPool, error) {
	if m.MarkCardAsDrawnFn != nil {
		return m.MarkCardAsDrawnFn(ctx, arg)
	}
	return sqlc.CardPool{}, nil
}

func (m *Querier) CollectCard(ctx context.Context, arg sqlc.CollectCardParams) error {
	if m.CollectCardFn != nil {
		return m.CollectCardFn(ctx, arg)
	}
	return nil
}

func (m *Querier) GetUserCollection(ctx context.Context, arg sqlc.GetUserCollectionParams) ([]sqlc.Card, error) {
	if m.GetUserCollectionFn != nil {
		return m.GetUserCollectionFn(ctx, arg)
	}
	return nil, nil
}

func (m *Querier) CardInCollection(ctx context.Context, arg sqlc.CardInCollectionParams) (bool, error) {
	if m.CardInCollectionFn != nil {
		return m.CardInCollectionFn(ctx, arg)
	}
	return false, nil
}
