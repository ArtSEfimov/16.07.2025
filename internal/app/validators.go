package app

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

const (
	jpg = ".jpg"
	pdf = ".pdf"
)

func validateURLFormat(request *LinkRequest) ([]Link, []Link) {
	newLinkValidator := validator.New()
	var validLinks = make([]Link, 0, len(request.Links))
	var invalidLinks = make([]Link, 0, len(request.Links))

	for _, link := range request.Links {
		err := newLinkValidator.Struct(link)
		if err != nil {
			invalidLinks = append(invalidLinks, link)
		} else {
			validLinks = append(validLinks, link)
		}
	}
	return validLinks, invalidLinks
}

func validateURLAccessible(links []Link) ([]Link, []Link) {
	validLinks := make([]Link, 0, len(links))
	invalidLinks := make([]Link, 0, len(links))
	for _, link := range links {
		if isURLAccessible(link.URL) {
			validLinks = append(validLinks, link)
		} else {
			invalidLinks = append(invalidLinks, link)
		}
	}
	return validLinks, invalidLinks
}

func isURLAccessible(urlString string) bool {
	response, responseErr := http.Head(urlString)
	if responseErr == nil && response.StatusCode == http.StatusOK {
		err := response.Body.Close()
		if err != nil {
			panic(err)
		}
		return true
	}

	request, requestCreateErr := http.NewRequest("GET", urlString, nil)
	if requestCreateErr != nil {
		panic(requestCreateErr)
	}
	request.Header.Set("Range", "bytes=0-0")

	response, responseErr = http.DefaultClient.Do(request)
	if responseErr != nil {
		return false
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {

		}
	}()

	return response.StatusCode == http.StatusOK || response.StatusCode == http.StatusPartialContent
}

func validateFileExtension(urlString string) (string, bool) {
	response, _ := http.Get(urlString)
	defer func() {
		err := response.Body.Close()
		if err != nil {
		}
	}()

	buf := make([]byte, 512)
	n, err := response.Body.Read(buf)
	if err != nil && err != io.EOF {
		panic(err)
	}

	contentType := http.DetectContentType(buf[:n])
	fmt.Println(contentType)
	switch contentType {
	case "application/pdf":
		return pdf, true
	case "image/jpeg", "image/jpg":
		return jpg, true
	default:
		return contentType, false
	}

}
