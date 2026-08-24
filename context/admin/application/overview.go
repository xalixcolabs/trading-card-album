package admin_application

import (
	"context"

	"com.xalixcolabs.trading-card-album/context/admin/model/dto"
	"com.xalixcolabs.trading-card-album/database"
)

func GetOverview(q database.Querier) (admin_dto.Overview, error) {
	ctx := context.Background()
	stats, err := q.GetOverviewStats(ctx)
	if err != nil {
		return admin_dto.Overview{}, err
	}
	return admin_dto.Overview{
		Albums:       stats.Albums,
		Users:        stats.Users,
		Cards:        stats.Cards,
		Participants: stats.Participants,
		Contacts:     stats.Contacts,
		Collected:    stats.Collected,
	}, nil
}