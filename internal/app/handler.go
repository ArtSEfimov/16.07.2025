package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go_test_task_4/pkg/response"
	"net/http"
	"time"
)

const (
	sessionCleanupInterval = 60 * time.Minute
	sessionID              = "session_id"
)

const (
	objectsLimit       = 3
	objectLimitMessage = "Object limit exceeded"
)

const getLinksPath = "/links"

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

	router.HandleFunc(fmt.Sprintf("POST %s", getLinksPath), handler.GetUserLinks())

}

func (handler *Handler) GetUserLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, cookieErr := r.Cookie(sessionID)
		var objects uint8
		if cookieErr != nil {
			uuidString := uuid.NewString()
			handler.objectsCounter[uuidString] = 0
			http.SetCookie(w, getNewCookie(uuidString))
		} else {
			sessionUUID := session.Value
			if objects, ok := handler.objectsCounter[sessionUUID]; !ok {
				handler.objectsCounter[sessionUUID] = 0
			} else {
				if objects >= objectsLimit {
					http.Error(w, objectLimitMessage, http.StatusBadRequest)
					return
				}
			}
		}

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
		if validLinks == nil {
			errorMessage := getErrorMessage(invalidLinks)
			linkResponse := LinkResponse{
				Result:       nil,
				ErrorMessage: errorMessage,
			}
			response.JsonResponse(w, &linkResponse, http.StatusBadRequest)
			return
		}
		if validLinks != nil {

		}
	}
}
