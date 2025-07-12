package app

import (
	"net/http"
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

	}
}
