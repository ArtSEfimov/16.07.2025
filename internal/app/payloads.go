package app

type Link struct {
	URL string `json:"url" validate:"url"`
}

type LinkRequest struct {
	Links []Link `json:"links"`
}

type Task struct {
	Id         uint64   `json:"id"`
	Status     string   `json:"status"`
	Links      []string `json:"links"`
	Errors     []string `json:"errors"`
	ArchiveURL string   `json:"archiveURL"`
}
