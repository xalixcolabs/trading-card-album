package contact_application

import (
	"context"

	"com.xalixcolabs.trading-card-album/context/contact/model/dto"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database"
)

func ListContacts(q database.Querier, user user_model.User) ([]contact_dto.Contact, error) {
	ctx := context.Background()
	rows, err := q.ListContacts(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	contacts := make([]contact_dto.Contact, 0, len(rows))
	for _, row := range rows {
		contacts = append(contacts, contact_dto.Contact{
			UserID:      row.ID,
			Name:        row.Name,
			Email:       row.Email,
			Github:      row.Github,
			Linkedin:    row.Linkedin,
			Web:         row.Web,
			Description: row.Description,
			ScannedAt:   row.ScannedAt,
		})
	}
	return contacts, nil
}