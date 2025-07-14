package app

type Link struct {
	URL           string `json:"url" validate:"url"`
	FileExtension string
}

type LinkRequest struct {
	Links []Link `json:"links"`
}

type Task struct {
	ID            uint64            `json:"id"`
	Status        string            `json:"status"`
	ValidLinks    []Link            `json:"valid_links"`
	InvalidLinks  []Link            `json:"invalid_links"`
	ErrorMessages map[string]string `json:"error_messages"`
	ArchiveURL    string            `json:"archiveURL"`
}
