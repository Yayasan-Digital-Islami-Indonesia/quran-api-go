package validator

// TODO: implement — see issue #12
//
// ValidateLang checks that lang is "id" or "en".
// Returns "id" as the default when lang is empty.
// Returns domain.ErrInvalidLang for any other value.
//
// Usage:
//   lang, err := validator.ValidateIDParam(c.Query("id"))
//   if err != nil {
//       response.BadRequest(c, "invalid param")
//       return
//   }

import (
	"math"
	"quran-api-go/internal/domain"
	"strconv"
)

func ValidateIDParam(id string) (string, error) {
	convertedId, err := strconv.Atoi(id)
	if err != nil {
		return "", domain.ErrInvalidIDParam
	}

	if math.IsNaN(float64(convertedId)) {
		return "", domain.ErrInvalidIDParam
	}

	if convertedId < 0 {
		return "", domain.ErrInvalidIDParam
	}

	return id, nil
}
