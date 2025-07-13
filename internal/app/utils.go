package app

import (
	"net/http"
)

func createErrorMessages(messages map[string]string, invalidLinks []Link, message string) {
	for _, invalidLink := range invalidLinks {
		if _, ok := messages[invalidLink.URL]; !ok {
			messages[invalidLink.URL] = message
		}
	}
}

func getNewCookie(uuid string) *http.Cookie {
	return &http.Cookie{
		Name:     sessionID,
		Value:    uuid,
		Path:     getLinksPath,
		HttpOnly: true,
		Secure:   false,
	}
}
