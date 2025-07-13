package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const sessionCleanupInterval = 60 * time.Minute

type Handler struct {
	objectsCounter map[string]uint8
}

func NewHandler(router *http.ServeMux) {
	handler := Handler{
		objectsCounter: make(map[string]uint8),
	}

	// session reset
	go func() {
		for range time.Tick(sessionCleanupInterval) {
			handler.objectsCounter = make(map[string]uint8)
		}
	}()

	router.HandleFunc("POST /links", handler.GetUserLinks())

}

func (handler *Handler) GetUserLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var links LinkRequest
		bodyReader := bufio.NewReader(r.Body)
		defer func() {
			err := r.Body.Close()
			if err != nil {
				panic(err)
			}
		}()

		bodyDecodeErr := json.NewDecoder(bodyReader).Decode(&links)
		if bodyDecodeErr != nil {
			http.Error(w, bodyDecodeErr.Error(), http.StatusBadRequest)
			return
		}
		validLinks, invalidLinks := ValidateURL(&links)
		if invalidLinks != nil {
			var invalidURLs strings.Builder
			for _, invalidLink := range invalidLinks {
				invalidURLs.WriteString(invalidLink.URL)
			}
			errorMessage := fmt.Sprintf("Invalid links: %s", invalidURLs.String())
			http.Error(w, errorMessage, http.StatusBadRequest)
		}
		if validLinks == nil {
			
		}
	}
}
