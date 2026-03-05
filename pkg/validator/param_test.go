package validator

import "testing"

func TestIOParamValidator(t *testing.T) {

func TestIDParamValidator(t *testing.T) {
	// Check if id is empty return error value
	id := ""
	id, err := ValidateIDParam(id)
	if err == nil {
		t.Logf("expect validation id param to return error if the id does not contain any")
		t.Fail()
	}

	// Check, Validate Param must return error if id is not a number
	id = "Get"
	id, err = ValidateIDParam(id)
	if err == nil {
		t.Logf("expect validation id param to return error if the id is not number")
		t.Fail()
	}

	// Check, Validate Param must return errror if id is negative value
	id = "-2"
	_, err = ValidateIDParam(id)
	if err == nil {
		t.Logf("expect validation id param to return error if the id is negative")
		t.Fail()
	}

	// Check, Validate Param must not return error if id is not negative value
	id = "1"
	id, err = ValidateIDParam(id)
	if err != nil {
		t.Logf("expect validation id param to accept id value")
		t.Fail()
	}
}
