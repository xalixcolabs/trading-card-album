package album_dto

type CreateAlbumRequest struct {
	Title string                   `json:"title"`
	Cards []CreateAlbumCardRequest `json:"cards"`
}

type CreateAlbumCardRequest struct {
	Number      string `json:"number"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
}
