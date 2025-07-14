package app

import (
	"fmt"
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

func (service *Service) CreateTask(userUUID string) {

}

func (service *Service) CreateFile(userUUID string) {

	//out, err := os.Create(filename)
	//if err != nil {
	//	return err
	//}
	//defer out.Close()
	//_, err = io.Copy(out, io.MultiReader(bytes.NewReader(buf[:n]), resp.Body))
	//if err != nil {
	//	return err
	//}
	//
	//fmt.Println("Сохранено как:", filename)
	//return nil
}
