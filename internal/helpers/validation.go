package helpers

import (
	"net/http"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// GetValidator returns a singleton validator instance with custom validations registered.
func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})
	return validate
}

// ValidationError represents a single field validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse is the response format for validation errors.
type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Details []ValidationError `json:"details"`
}

// ValidateStruct validates a struct and returns a formatted error response if invalid.
// returns nil if validation passes.
func ValidateStruct(s interface{}) *ValidationErrorResponse {
	v := GetValidator()
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return &ValidationErrorResponse{
			Error: "Validation failed",
			Details: []ValidationError{
				{Field: "unknown", Message: err.Error()},
			},
		}
	}

	var details []ValidationError
	for _, e := range validationErrors {
		details = append(details, ValidationError{
			Field:   toSnakeCase(e.Field()),
			Message: formatValidationMessage(e),
		})
	}

	return &ValidationErrorResponse{
		Error:   "Validation failed",
		Details: details,
	}
}

// BindAndValidate binds the request body to the given struct and validates it.
// returns an error response to send to the client if binding or validation fails.
func BindAndValidate(c echo.Context, s interface{}) error {
	if err := c.Bind(s); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if validationErr := ValidateStruct(s); validationErr != nil {
		return c.JSON(http.StatusBadRequest, validationErr)
	}

	return nil
}

// formatValidationMessage creates a human-readable validation message.
func formatValidationMessage(e validator.FieldError) string {
	field := toSnakeCase(e.Field())

	switch e.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + e.Param() + " characters"
	case "max":
		return field + " must be at most " + e.Param() + " characters"
	case "url":
		return field + " must be a valid URL"
	case "gt":
		return field + " must be greater than " + e.Param()
	case "gte":
		return field + " must be greater than or equal to " + e.Param()
	case "lt":
		return field + " must be less than " + e.Param()
	case "lte":
		return field + " must be less than or equal to " + e.Param()
	default:
		return field + " failed validation: " + e.Tag()
	}
}

// toSnakeCase converts PascalCase or camelCase to snake_case.
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
