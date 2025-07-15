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

	router.HandleFunc(fmt.Sprintf("POST %s", createTaskPath), handler.CreateTask())
	router.HandleFunc(fmt.Sprintf("POST %s", addLinkPath), handler.AddLink())
	router.HandleFunc(fmt.Sprintf("GET %s", getTaskStatusPath), handler.GetTaskStatus())

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
			
			http.Error(w, bodyDecodeErr.Error(), http.StatusBadRequest)
			return
		} 
		task := Task{
			ID:            handler.service.GetTaskID(),
			Status:        taskStatusCreated,
			ValidLinks:    make([]Link, 0),
			InvalidLinks:  make([]Link, 0),
			ErrorMessages: make(map[string]string),
			ArchiveURL:    "",
		}
		handler.repository.AddUserTask(userUUID, task)
			

		// 
		
		if links.Links == nil {
			response.JsonResponse(w, &task, http.StatusCreated)
			return
		}

		validLinks, invalidLinks := validateURLFormat(&links)
		if validLinks == nil {
			task.InvalidLinks = append(task.InvalidLinks, invalidLinks...)
			addErrorMessages(task.ErrorMessages, invalidLinks, errInvalidLinkFormat)

			response.JsonResponse(w, &task, http.StatusCreated)
			return
		}

		if validLinks != nil {
			validLinks, invalidLinks := validateURLAccessible(validLinks)
			task.InvalidLinks = append(task.InvalidLinks, invalidLinks...)
			addErrorMessages(task.ErrorMessages, invalidLinks, errInaccessibleLink)
			
			for _, link := range validLinks {
				ext, isValid := validateFileExtension(link.URL)
					
				enrichedLink := Link{
					URL:           link.URL,
					FileExtension: ext,
				}
				if isValid{
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
					task.InvalidLinks = append(task.InvalidLinks, enrichedLink) 
					task.ErrorMessages[enrichedLink.URL] = fmt.Sprintf("%s: %s", errUnsupportedContentType, ext)
					}
				
			

				]
	
			response.JsonResponse(w, &task, http.StatusCreated)
		
	}
}


func (handler *Handler) GetTaskStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, cookieErr := r.Cookie(sessionID)
		userUUID = session.Value
		if handler.repository.isUserHasTask(userUUID) {
				taskID := handler.repository.GetUserTaskID(userUUID)
				task = handler.repository.GetTaskByID(taskID)
			response.JsonResponse(w, &task, http.StatusOK)
			return
			}
		
		http.Error(w, errUserHasNoTasks, http.StatusBadRequest)
			return
	}
}


func (handler *Handler) AddLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, cookieErr := r.Cookie(sessionID)
		userUUID = session.Value
		if !handler.repository.isUserHasTask(userUUID) {
			
			http.Error(w, errUserHasNoTasks, http.StatusBadRequest)
			return
			
		}
		taskID := handler.repository.GetUserTaskID(userUUID)
		task = handler.repository.GetTaskByID(taskID)
		
		var links LinkRequest
		bodyDecodeErr := json.NewDecoder(bodyReader).Decode(&links)
		if bodyDecodeErr != nil && bodyDecodeErr != io.EOF { // invalid json
			http.Error(w, bodyDecodeErr.Error(), http.StatusBadRequest)
			return
		}

		if links.Links == nil {
			response.JsonResponse(w, &task, http.StatusCreated)
			return
		}

		validLinks, invalidLinks := validateURLFormat(&links)
		if validLinks == nil {
			task.InvalidLinks = append(task.InvalidLinks, invalidLinks...)
			addErrorMessages(task.ErrorMessages, invalidLinks, errInvalidLinkFormat)

			response.JsonResponse(w, &task, http.StatusOK)
			return
		}
		if validLinks != nil {
			validLinks, invalidLinks := validateURLAccessible(validLinks)
			task.InvalidLinks = append(task.InvalidLinks, invalidLinks...)
			addErrorMessages(task.ErrorMessages, invalidLinks, errInaccessibleLink)
			
			for _, link := range validLinks {
				ext, isValid := validateFileExtension(link.URL)
					
				enrichedLink := Link{
					URL:           link.URL,
					FileExtension: ext,
				}
				if isValid{
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
					task.InvalidLinks = append(task.InvalidLinks, enrichedLink) 
					task.ErrorMessages[enrichedLink.URL] = fmt.Sprintf("%s: %s", errUnsupportedContentType, ext)
					}
				
			

				]
	
			response.JsonResponse(w, &task, http.StatusCreated)
		
	}
}
