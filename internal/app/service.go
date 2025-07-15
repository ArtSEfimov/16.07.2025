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


func CreateZipFile(files []string) error {
	const outputZipPath = "tmp/"
	outputZipFile, zipCreationErr := os.Create(outputZipPath)
	if zipCreationErr != nil {
		return zipCreationErr
	}
	defer outputZipFile.Close()


	zipWriter := zip.NewWriter(outputZipFile)
	defer zipWriter.Close()

	for i, url := range urls {
		
    		resp, err := http.Get(url)
    		if err != nil {
        		return fmt.Errorf("failed to download %s: %w", url, err)
   		}
    	defer resp.Body.Close()

    
        fileName = fmt.Sprintf("file%d%s", i, ext) // fallback


   	 header := &zip.FileHeader{
        	Name:   fileName,
        	Method: zip.Deflate,
    	}

    	writer, err := zipWriter.CreateHeader(header)
    	if err != nil {
	        return err
    	}

    	if _, err := io.Copy(writer, resp.Body); err != nil {
        	return fmt.Errorf("failed to copy data from %s: %w", url, err)
    	}

	return nil
}


func (service *Service) CreateTask(userUUID string, task Task) {
	
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
