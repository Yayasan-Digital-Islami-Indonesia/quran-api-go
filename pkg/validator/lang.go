package validator

// TODO: implement — see issue #12
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
	_ = domain.ErrInvalidLang // sentinel used as return value
	panic("not implemented")
}
