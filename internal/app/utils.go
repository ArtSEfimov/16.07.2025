package app

import (
	"net/http"
)

func createErrorMessages(messages map[string]string, invalidLinks map[Link]struct{}, message string) {
	for invalidLink := range invalidLinks {
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

func createSlice(linkMap map[Link]struct{}) []Link {
	linkSlice := make([]Link, 0, len(linkMap))
	for link := range linkMap {
		linkSlice = append(linkSlice, link)
	}
	return linkSlice
}
