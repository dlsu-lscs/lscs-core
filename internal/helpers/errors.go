package helpers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorResponse represents a standardized error response for Swagger documentation
type ErrorResponse struct {
	Error   string `json:"error" example:"Something went wrong"`
	Code    string `json:"code,omitempty" example:"BAD_REQUEST"`
	Details any    `json:"details,omitempty"`
}

// APIError represents a standardized error response.
type APIError struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details any    `json:"details,omitempty"`
}

// common error codes
const (
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeInternal       = "INTERNAL_ERROR"
	ErrCodeBadRequest     = "BAD_REQUEST"
	ErrCodeNotImplemented = "NOT_IMPLEMENTED"
)

// NewAPIError creates a new API error with the given message and optional code.
func NewAPIError(message string, code string) *APIError {
	return &APIError{
		Error: message,
		Code:  code,
	}
}

// WithDetails adds details to the error.
func (e *APIError) WithDetails(details any) *APIError {
	e.Details = details
	return e
}

// --- Error response helpers ---

// ErrBadRequest returns a 400 Bad Request response.
func ErrBadRequest(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, NewAPIError(message, ErrCodeBadRequest))
}

// ErrUnauthorized returns a 401 Unauthorized response.
func ErrUnauthorized(c echo.Context, message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return c.JSON(http.StatusUnauthorized, NewAPIError(message, ErrCodeUnauthorized))
}

// ErrForbidden returns a 403 Forbidden response.
func ErrForbidden(c echo.Context, message string) error {
	if message == "" {
		message = "Forbidden"
	}
	return c.JSON(http.StatusForbidden, NewAPIError(message, ErrCodeForbidden))
}

// ErrNotFound returns a 404 Not Found response.
func ErrNotFound(c echo.Context, message string) error {
	if message == "" {
		message = "Not found"
	}
	return c.JSON(http.StatusNotFound, NewAPIError(message, ErrCodeNotFound))
}

// ErrConflict returns a 409 Conflict response.
func ErrConflict(c echo.Context, message string) error {
	return c.JSON(http.StatusConflict, NewAPIError(message, ErrCodeConflict))
}

// ErrInternal returns a 500 Internal Server Error response.
func ErrInternal(c echo.Context, message string) error {
	if message == "" {
		message = "Internal server error"
	}
	return c.JSON(http.StatusInternalServerError, NewAPIError(message, ErrCodeInternal))
}

// ErrValidation returns a 400 Bad Request with validation details.
func ErrValidation(c echo.Context, details *ValidationErrorResponse) error {
	return c.JSON(http.StatusBadRequest, details)
}

// --- Specific error responses used across the app ---

// ErrNotMember returns a standardized "not an LSCS member" response.
func ErrNotMember(c echo.Context, identifier string, identifierType string) error {
	resp := map[string]string{
		"error":        "Not an LSCS member",
		"state":        "absent",
		identifierType: identifier,
	}
	return c.JSON(http.StatusNotFound, resp)
}
