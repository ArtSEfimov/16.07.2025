package app

type Link struct {
	URL string `json:"url" validate:"required,url"`
}
type LinkRequest struct {
	Links []Link `json:"links"`
}
