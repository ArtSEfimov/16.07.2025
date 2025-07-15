package app

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
)

func addErrorMessages(messages map[string]string, invalidLinks []Link, message string) {
	for _, invalidLink := range invalidLinks {
		if _, ok := messages[invalidLink.URL]; !ok {
			messages[invalidLink.URL] = message
		}
	}
}

func getNewCookie(uuid string) *http.Cookie {
	return &http.Cookie{
		Name:     sessionID,
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	}
}

func createZipFile(task *Task, errChan chan error, outputZipPath string) {
	defer close(errChan)

	if err := os.MkdirAll(baseOutputZipPath, 0755); err != nil {
		errChan <- fmt.Errorf("failed to create base dir %q: %w", baseOutputZipPath, err)
		return
	}

	outputZipFile, err := os.Create(outputZipPath)
	if err != nil {
		errChan <- fmt.Errorf("failed to create zip file: %w", err)
		return
	}
	defer outputZipFile.Close()

	zipWriter := zip.NewWriter(outputZipFile)
	defer zipWriter.Close()

	for i, link := range task.ValidLinks {
		url := link.URL

		response, err := http.Get(url)
		if err != nil {
			errChan <- fmt.Errorf("failed to download %s: %w", url, err)
			return
		}

		func() {
			defer response.Body.Close()

			fileName := fmt.Sprintf("file%d%s", i, link.FileExtension)

			header := &zip.FileHeader{
				Name:   fileName,
				Method: zip.Deflate,
			}

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				errChan <- fmt.Errorf("failed to create zip writer header: %w", err)
				return
			}

			if _, err = io.Copy(writer, response.Body); err != nil {
				errChan <- fmt.Errorf("failed to copy data from %s: %w", url, err)
				return
			}
		}()
	}

	errChan <- nil
}
