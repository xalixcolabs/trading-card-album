package card_dto

type CreateCardRequest struct {
	AlbumId string `json:"album_id"`
	Number string `json:"number"`
	Name string `json:"name"`
	Description string `json:"description"`
	ImageUrl string `json:"image_url"`
}