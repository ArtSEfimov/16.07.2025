package app

import (
	"github.com/go-playground/validator/v10"
	"net/url"
	"path"
	"strings"
)

const (
	jpg  = ".jpg"
	jpeg = ".jpeg"
	pdf  = ".pdf"
)

var validExtensions = map[string]struct{}{jpg: {}, jpeg: {}, pdf: {}}

func validateURL(request *LinkRequest) ([]Link, []Link) {
	newLinkValidator := validator.New()
	var validLinks []Link
	var invalidLinks []Link

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

func validateObjectExtension(urlString string) (string, bool) {
	urlObject, _ := url.Parse(urlString)
	urlObjectPath := urlObject.Path
	ext := strings.ToLower(path.Ext(urlObjectPath))
	if _, ok := validExtensions[ext]; ok {
		return ext, true
	}
	return "", false
}
