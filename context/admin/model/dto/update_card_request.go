package admin_dto

type UpdateCardRequest struct {
	Number      string `json:"number"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
}