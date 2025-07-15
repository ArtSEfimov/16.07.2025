package app

import (
	"fmt"
	"path/filepath"
)

type Service struct {
	taskID uint64
}

func NewService() *Service {
	return &Service{
		taskID: 1,
	}
}

func (service *Service) GetTaskID() uint64 {
	taskID := service.taskID
	service.taskID++
	return taskID
}

func (service *Service) DownloadFile(objectURL string) error {
	objectExt, isValid := validateFileExtension(objectURL)
	if !isValid {
		return fmt.Errorf("invalid file extension %service", objectExt)
	}
	return nil
}

func (service *Service) CreateZipFile(task *Task) {
	errChan := make(chan error)

	filename := fmt.Sprintf("Archive%d.zip", task.ID)
	outputZipPath := filepath.Join(baseOutputZipPath, filename)

	task.Status = taskStatusProcessing

	go createZipFile(task, errChan, outputZipPath)

	go func() {
		err := <-errChan
		if err != nil {
			task.Status = taskStatusError
			task.ErrorMessages[fmt.Sprintf("task %d %s", task.ID, errZipFileCreation)] = err.Error()
		} else {
			task.Status = taskStatusCompleted
			task.ArchiveURL = fmt.Sprintf("http://localhost:8080/%s/%s", baseOutputZipPath, filename)
		}
	}()
}
