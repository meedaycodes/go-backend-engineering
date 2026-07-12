// Package validator provides request validation functions that check input
// format before it reaches the service layer. Collects all validation errors
// into a map (field name → message) so clients see every problem at once,
// not one at a time.
package validator

import (
	"strings"

	"github.com/meedaycodes/day14-integration-testing/internal/model"
)

// ValidationError holds field-level validation failures. Satisfies the error
// interface via Error(). The Fields map is encoded as JSON in the handler,
// giving clients a structured response they can parse per-field.
type ValidationError struct {
	Fields map[string]string
}

func (v *ValidationError) Error() string {
	return "Validation failed"
}

// ValidateCreateUser checks format constraints on signup input: name required
// and ≤100 chars, email required and contains @, password required and ≥8 chars.
// Uses else-if chains so empty fields don't trigger format checks (e.g., empty
// email won't also fail the @ check).
func ValidateCreateUser(req model.CreateUserRequest) error {

	errorMap := make(map[string]string)

	if req.Name == "" {
		errorMap["name"] = "name is required"
	}
	if len(req.Name) > 100 {
		errorMap["name"] = "name must be 100 characters or less"
	}
	if req.Email == "" {
		errorMap["email"] = "email is required"
	} else if !strings.Contains(req.Email, "@") {
		errorMap["email"] = "email format is invalid"
	}
	if req.Password == "" {
		errorMap["password"] = "password is required"
	} else if len(req.Password) < 8 {
		errorMap["password"] = "password must be at least 8 characters"
	}

	if len(errorMap) > 0 {
		return &ValidationError{Fields: errorMap}
	}

	return nil

}

// ValidateLoginRequest checks that email and password are present.
// No format checks — login just needs non-empty credentials.
func ValidateLoginRequest(req model.LoginRequest) error {

	errorMap := make(map[string]string)
	if req.Email == "" {
		errorMap["email"] = "email is required"
	}
	if req.Password == "" {
		errorMap["password"] = "password is required"
	}
	if len(errorMap) > 0 {
		return &ValidationError{Fields: errorMap}
	}
	return nil
}
