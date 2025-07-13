package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go_test_task_4/pkg/response"
	"io"
	"net/http"
)

type Handler struct {
	service    *Service
	repository *Repository
}

func NewHandler(router *http.ServeMux) {
	handler := Handler{}

	router.HandleFunc(fmt.Sprintf("POST %s", getLinksPath), handler.CreateTask())

}

func (handler *Handler) CreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if handler.repository.GetTaskCount() >= taskLimit {
			http.Error(w, taskLimitMessage, http.StatusTooManyRequests)
			return
		}

		session, cookieErr := r.Cookie(sessionID)
		var userUUID string
		if cookieErr != nil {
			userUUID = uuid.NewString()
			handler.repository.AddNewUser(userUUID)
			http.SetCookie(w, getNewCookie(userUUID))
		} else {
			userUUID = session.Value
		}

		bodyReader := bufio.NewReader(r.Body)
		defer func() {
			err := r.Body.Close()
			if err != nil {
				panic(err)
			}
		}()

		var links LinkRequest
		bodyDecodeErr := json.NewDecoder(bodyReader).Decode(&links)
		if bodyDecodeErr != nil && bodyDecodeErr != io.EOF { // invalid json
			http.Error(w, bodyDecodeErr.Error(), http.StatusBadRequest)
			return
		}
		if links.Links == nil {
			linkResponse := Task{
				ID:            handler.service.GetTaskID(),
				Status:        taskStatusCreated,
				ValidLinks:    nil,
				InvalidLinks:  nil,
				ErrorMessages: nil,
				ArchiveURL:    "",
			}
			response.JsonResponse(w, &linkResponse, http.StatusCreated)
			return
		}

		validLinks, invalidLinks := validateURLFormat(&links)
		var errorMessages map[string]string
		if validLinks == nil {
			createErrorMessages(errorMessages, invalidLinks, errInvalidLinkFormat)
			linkResponse := Task{
				ID:            handler.service.GetTaskID(),
				Status:        taskStatusCreated,
				ValidLinks:    nil,
				InvalidLinks:  invalidLinks,
				ErrorMessages: errorMessages,
				ArchiveURL:    "",
			}
			response.JsonResponse(w, &linkResponse, http.StatusCreated)
			return
		}
		if validLinks != nil {
			for _, link := range validLinks {
				if isURLAccessible(link.URL) {
					if ext, isValidExt := validateObjectExtension(link.URL); isValidExt {
						link.ObjectExtension = ext
						handler.repository.AddNewObject(userUUID, link)
						// TODO счиатем количество ссылок для юзера и действуем соответсвенно
					} else {
						invalidLinks = append(invalidLinks, link)
						errorMessages[link.URL] = errUnsupportedFileType
					}
				} else {
					invalidLinks = append(invalidLinks, link)
					errorMessages[link.URL] = errInaccessibleLink
				}

			}
		}

	}
}
