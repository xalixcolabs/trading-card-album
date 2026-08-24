package album_participant_dto

type ShareAssignedCardResponse struct {
	ContactId string `json:"contact_id"`
	AlbumId   string `json:"album_id"`
	CardId    string `json:"card_id"`
	Secret    string `json:"secret"`
}
