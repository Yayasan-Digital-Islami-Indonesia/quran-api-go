package validator

//
// ValidateLang checks that lang is "id" or "en".
// Returns "id" as the default when lang is empty.
// Returns domain.ErrInvalidLang for any other value.
//
// Usage:
//   lang, err := validator.ValidateLang(c.Query("lang"))
//   if err != nil {
//       response.BadRequest(c, "lang must be 'id' or 'en'")
//       return
//   }

import "quran-api-go/internal/domain"

func ValidateLang(lang string) (string, error) {
	if lang == "" {
		return "id", nil
	}

	if lang != "id" && lang != "en" {
		return "", domain.ErrInvalidLang
	}

	return lang, nil
}
