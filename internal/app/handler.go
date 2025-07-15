package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go_test_task_4/pkg/response"
	"io"
	"net/http"
	"strconv"
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

	router.HandleFunc(fmt.Sprintf("POST %s", createTaskPath), handler.CreateTask())
	router.HandleFunc(fmt.Sprintf("POST %s/{id}", addLinkPath), handler.AddLink())
	router.HandleFunc(fmt.Sprintf("GET %s/{id}", getTaskStatusPath), handler.GetTaskStatus())

}

func (handler *Handler) CreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if handler.repository.GetActiveTaskCount() >= taskLimit {
			http.Error(w, taskLimitMessage, http.StatusTooManyRequests)
			return
		}

		session, cookieErr := r.Cookie(sessionID)
		var userUUID string
		if cookieErr != nil {
			userUUID = uuid.NewString()
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
		var task *Task
		bodyDecodeErr := json.NewDecoder(bodyReader).Decode(&links)
		if bodyDecodeErr != nil && bodyDecodeErr != io.EOF { // invalid json
			http.Error(w, bodyDecodeErr.Error(), http.StatusBadRequest)
			return
		}

		task = &Task{
			ID:            handler.service.GetTaskID(),
			Status:        taskStatusCreated,
			ValidLinks:    make([]Link, 0),
			InvalidLinks:  make([]Link, 0),
			ErrorMessages: make(map[string]string),
			ArchiveURL:    "",
		}
		handler.repository.AddTask(userUUID, task)

		if links.Links == nil {
			response.JsonResponse(w, task, http.StatusCreated)
			return
		}

		validLinks, invalidLinks := validateURLFormat(&links)
		if validLinks == nil {
			task.InvalidLinks = append(task.InvalidLinks, invalidLinks...)
			addErrorMessages(task.ErrorMessages, invalidLinks, errInvalidLinkFormat)
			response.JsonResponse(w, task, http.StatusCreated)
			return
		}

		if validLinks != nil {
			validLinks, invalidLinks = validateURLAccessible(validLinks)
			task.InvalidLinks = append(task.InvalidLinks, invalidLinks...)
			addErrorMessages(task.ErrorMessages, invalidLinks, errInaccessibleLink)

			for _, link := range validLinks {
				ext, isValid := validateFileExtension(link.URL)

				enrichedLink := Link{
					URL:           link.URL,
					FileExtension: ext,
				}
				if isValid {
					if task.getValidLinksCount() < filesLimit {
						task.ValidLinks = append(task.ValidLinks, enrichedLink)
						task.Status = taskStatusPending
					}
					if task.getValidLinksCount() == filesLimit && (task.Status == taskStatusCreated || task.Status == taskStatusPending) {
						handler.service.CreateZipFile(task)
					}
				} else {
					task.InvalidLinks = append(task.InvalidLinks, enrichedLink)
					task.ErrorMessages[enrichedLink.URL] = fmt.Sprintf("%s: %s", errUnsupportedContentType, ext)
				}

			}
		}
		response.JsonResponse(w, task, http.StatusCreated)
	}
}

func (handler *Handler) GetTaskStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, cookieErr := r.Cookie(sessionID)
		if cookieErr != nil {
			http.Error(w, errUserNotFound, http.StatusBadRequest)
			return
		}
		userUUID := session.Value

		idString := r.PathValue("id")
		id, parseErr := strconv.ParseUint(idString, 10, 64)
		if parseErr != nil {
			idErr := fmt.Errorf("wrong id format: %w, get %s", parseErr, idString)
			http.Error(w, idErr.Error(), http.StatusBadRequest)
			return
		}

		task := handler.repository.GetTaskByID(userUUID, id)
		if task == nil {
			http.Error(w, fmt.Sprintf("%s %d", errUserHasNoTaskByID, id), http.StatusBadRequest)
			return
		}

		response.JsonResponse(w, task, http.StatusOK)
	}
}

func (handler *Handler) AddLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, cookieErr := r.Cookie(sessionID)
		if cookieErr != nil {
			http.Error(w, errUserNotFound, http.StatusBadRequest)
			return
		}
		userUUID := session.Value
		if !handler.repository.isUserHasTask(userUUID) {
			http.Error(w, errUserHasNoTaskByID, http.StatusBadRequest)
			return
		}

		idString := r.PathValue("id")
		id, parseErr := strconv.ParseUint(idString, 10, 64)
		if parseErr != nil {
			idErr := fmt.Errorf("wrong id format: %w, get %s", parseErr, idString)
			http.Error(w, idErr.Error(), http.StatusBadRequest)
			return
		}

		task := handler.repository.GetTaskByID(userUUID, id)
		if task == nil {
			http.Error(w, fmt.Sprintf("%s %d", errUserHasNoTaskByID, id), http.StatusBadRequest)
			return
		}
		if task.Status == taskStatusCompleted || task.Status == taskStatusError || task.Status == taskStatusProcessing {
			http.Error(w, fmt.Sprintf("%s task status is %s", errCannotAddLinkToTask, task.Status), http.StatusBadRequest)
			return
		}

		bodyReader := bufio.NewReader(r.Body)
		defer func() {
			err := r.Body.Close()
			if err != nil {
				panic(err)
			}
		}()

		// для устранения дублирования кода в методе CreateTask нужно передавать в функцию много зависимостей
		// поэтому пока что избыточность сохранится

		var links LinkRequest
		bodyDecodeErr := json.NewDecoder(bodyReader).Decode(&links)
		if bodyDecodeErr != nil && bodyDecodeErr != io.EOF { // invalid json
			http.Error(w, bodyDecodeErr.Error(), http.StatusBadRequest)
			return
		}

		if links.Links == nil {
			response.JsonResponse(w, task, http.StatusCreated)
			return
		}

		validLinks, invalidLinks := validateURLFormat(&links)
		if validLinks == nil {
			task.InvalidLinks = append(task.InvalidLinks, invalidLinks...)
			addErrorMessages(task.ErrorMessages, invalidLinks, errInvalidLinkFormat)
			response.JsonResponse(w, task, http.StatusOK)
			return
		}

		if validLinks != nil {
			validLinks, invalidLinks = validateURLAccessible(validLinks)
			task.InvalidLinks = append(task.InvalidLinks, invalidLinks...)
			addErrorMessages(task.ErrorMessages, invalidLinks, errInaccessibleLink)
			for _, link := range validLinks {
				ext, isValid := validateFileExtension(link.URL)

				enrichedLink := Link{
					URL:           link.URL,
					FileExtension: ext,
				}

				if isValid {
					if task.getValidLinksCount() < filesLimit {
						task.ValidLinks = append(task.ValidLinks, enrichedLink)
						task.Status = taskStatusPending
					}
					if task.getValidLinksCount() == filesLimit && (task.Status == taskStatusCreated || task.Status == taskStatusPending) {
						handler.service.CreateZipFile(task)
					}
				} else {
					task.InvalidLinks = append(task.InvalidLinks, enrichedLink)
					task.ErrorMessages[enrichedLink.URL] = fmt.Sprintf("%s: %s", errUnsupportedContentType, ext)
				}
			}
		}
		response.JsonResponse(w, task, http.StatusOK)
	}
}
