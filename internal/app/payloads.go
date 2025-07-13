package app

type Link struct {
	URL string `json:"url" validate:"required,url"`
}

type LinkRequest struct {
	Links []Link `json:"links"`
}

type LinkResponse struct {
	Result       any
	ErrorMessage string `json:"error_message"`
}
