package app

import (
	"fmt"
	"net/http"
	"strings"
)

func getErrorMessage(invalidLinks []Link) string {
	var invalidURLs strings.Builder
	for _, invalidLink := range invalidLinks {
		invalidURLs.WriteString(invalidLink.URL)
	}

	return fmt.Sprintf("Invalid links: %s", invalidURLs.String())
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
