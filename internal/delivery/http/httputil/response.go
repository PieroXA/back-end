package httputil

import (
	"encoding/json"
	"errors"
	"net/http"

	"mifinca-backend/internal/domain"
)

// APIResponse is the uniform JSON envelope for all API responses.
type APIResponse struct {
	Data  any       `json:"data"`
	Error *APIError `json:"error"`
}

// APIError is the structured error payload.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// JSON writes a success response with the given status code and payload.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse{Data: data, Error: nil})
}

// Error maps domain errors to the correct HTTP status code and writes the error response.
func Error(w http.ResponseWriter, err error) {
	code, apiErr := mapError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(APIResponse{Data: nil, Error: apiErr})
}

// ErrorWithMessage writes a custom error response.
func ErrorWithMessage(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse{
		Data:  nil,
		Error: &APIError{Code: code, Message: message},
	})
}

func mapError(err error) (int, *APIError) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, &APIError{Code: "NOT_FOUND", Message: err.Error()}
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict, &APIError{Code: "CONFLICT", Message: err.Error()}
	case errors.Is(err, domain.ErrValidation):
		return http.StatusBadRequest, &APIError{Code: "VALIDATION", Message: err.Error()}
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, &APIError{Code: "FORBIDDEN", Message: err.Error()}
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized, &APIError{Code: "UNAUTHORIZED", Message: err.Error()}
	default:
		return http.StatusInternalServerError, &APIError{Code: "INTERNAL", Message: "internal server error"}
	}
}
