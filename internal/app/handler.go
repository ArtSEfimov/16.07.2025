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

func NewHandler(router *http.ServeMux, repository *Repository, service *Service) {
	handler := Handler{
		service:    service,
		repository: repository,
	}

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
		var task Task
		bodyDecodeErr := json.NewDecoder(bodyReader).Decode(&links)
		if bodyDecodeErr != nil && bodyDecodeErr != io.EOF { // invalid json
			http.Error(w, bodyDecodeErr.Error(), http.StatusBadRequest)
			return
		}
		
		// Check or create Task
	
		if handler.repository.isUserHasTask(userUUID) {
				taskID := handler.repository.GetUserTaskID(userUUID)
				task = handler.repository.GetTaskByID(taskID)
		} else {
				task := Task{
					ID:            handler.service.GetTaskID(),
					Status:        taskStatusCreated,
					ValidLinks:    make([]Link, 0),
					InvalidLinks:  make([]Link, 0),
					ErrorMessages: nil,
					ArchiveURL:    "",
				}
				handler.repository.AddUserTask(userUUID, task)
			}

		// 
		
		if links.Links == nil {
			
			response.JsonResponse(w, &task, http.StatusCreated)
			return
		}

		validLinks, invalidLinks := validateURLFormat(&links)
		var errorMessages = make(map[string]string)
		if validLinks == nil {
			createErrorMessages(errorMessages, invalidLinks, errInvalidLinkFormat)
			if handler.repository.isUserHasTask(userUUID) {
				taskID := handler.repository.GetUserTaskID(userUUID)
				task = handler.repository.GetTaskByID(taskID)
				task.ErrorMessages = errorMessages
				task.InvalidLinks = append(task.InvalidLinks, invalidLinks...)
			} else {
				task := Task{
					ID:            handler.service.GetTaskID(),
					Status:        taskStatusCreated,
					ValidLinks:    nil,
					InvalidLinks:  invalidLinksSlice,
					ErrorMessages: errorMessages,
					ArchiveURL:    "",
				}
				
				handler.repository.AddUserTask(userUUID, task)
			}
			
			response.JsonResponse(w, &task, http.StatusCreated)
			return
		}

		if validLinks != nil {
			if handler.repository.isUserHasTask(userUUID) {
				taskID := handler.repository.GetUserTaskID(userUUID)
				task = handler.repository.GetTaskByID(taskID)
			} else {

				task = Task{
					ID:           handler.service.GetTaskID(),
					Status:       taskStatusCreated,
					ValidLinks:   make([]Link, 0, len(validLinks)),
					InvalidLinks: make([]Link, 0, len(invalidLinks)),
				}

				handler.repository.AddUserTask(userUUID, task)
			}
			for link := range validLinks {
				var enrichedLink Link
				if isURLAccessible(link.URL) {
					if ext, isValidExt := validateFileExtension(link.URL); isValidExt {
						enrichedLink = Link{
							URL:           link.URL,
							FileExtension: ext,
						}
						task.ValidLinks = append(task.ValidLinks, enrichedLink)
						handler.repository.AddNewUserLink(userUUID, enrichedLink)
						if handler.repository.GetUserLinksCount(userUUID) == filesLimit {
							//task.Status = taskStatusProcessing
							// TODO запускаем сервис по созданию архива
							fmt.Println("Service IS RUNNING")

						} else if handler.repository.GetUserLinksCount(userUUID) > filesLimit {
							break
						} else {
							task.Status = taskStatusPending
						}
					} else {
						enrichedLink = Link{
							URL:           link.URL,
							FileExtension: ext,
						}
						invalidLinks[enrichedLink] = struct{}{}
						errorMessages[enrichedLink.URL] = errUnsupportedContentType
					}
				} else {
					invalidLinks[link] = struct{}{}
					errorMessages[link.URL] = errInaccessibleLink
				}
			}

			invalidLinksSlice := createSlice(invalidLinks)
			task.InvalidLinks = invalidLinksSlice
			task.ErrorMessages = errorMessages

			fmt.Println(task.ValidLinks)
			fmt.Println(task.InvalidLinks)

			response.JsonResponse(w, &task, http.StatusCreated)
		}
	}
}
