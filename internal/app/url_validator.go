package app

import (
	"github.com/go-playground/validator/v10"
)

func ValidateURL(request *LinkRequest) ([]Link, []Link) {
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
