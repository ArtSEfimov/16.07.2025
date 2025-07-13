package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go_test_task_4/pkg/response"
	"net/http"
)

const sessionID = "session_id"

const (
	taskLimit        = 3
	taskLimitMessage = "task limit exceeded"
)

const getLinksPath = "/links"

type Handler struct {
	service    *Service
	repository *Repository
}

func NewHandler(router *http.ServeMux) {
	handler := Handler{}

	router.HandleFunc(fmt.Sprintf("POST %s", getLinksPath), handler.GetUserLinks())

}

func (handler *Handler) CreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, cookieErr := r.Cookie(sessionID)
		if cookieErr != nil {
			uuidString := uuid.NewString()
			handler.repository.AddNewUser(uuidString)
			http.SetCookie(w, getNewCookie(uuidString))
		} else {
			sessionUUID := session.Value
			userTaskCount := handler.repository.GetUserTaskCount(sessionUUID)
			if userTaskCount >= taskLimit {
				http.Error(w, taskLimitMessage, http.StatusTooManyRequests)
				return
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
		validLinks, invalidLinks := validateURL(&links)
		if validLinks == nil {
			//errorMessage := getErrorMessage(invalidLinks)
			//linkResponse := Task{
			//	Result:       nil,
			//	ErrorMessage: errorMessage,
			//}
			response.JsonResponse(w, &linkResponse, http.StatusBadRequest)
			return
		}
		if validLinks != nil {
			for counter := objects; counter <= 3; counter++ {

			}
		}
	}
}
