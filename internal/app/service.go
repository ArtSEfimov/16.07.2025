package app

import "fmt"

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) DownloadObject(objectURL string) error {
	objectExt, isValid := validateObjectExtension(objectURL)
	if !isValid {
		return fmt.Errorf("invalid file extension %s", objectExt)
	}
	return nil

}
