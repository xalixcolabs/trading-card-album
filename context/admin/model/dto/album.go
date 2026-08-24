package admin_dto

type Album struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	TotalCards       int64  `json:"total_cards"`
	CreatedAt        int64  `json:"created_at"`
	ParticipantCount int64  `json:"participant_count"`
}