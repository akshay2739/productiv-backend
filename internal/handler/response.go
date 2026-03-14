package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/akshay/productiv-backend/internal/domain"
)

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// ErrorResponse represents an error response body.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// HandleError maps domain errors to appropriate HTTP responses.
func HandleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		JSON(w, http.StatusNotFound, ErrorResponse{Error: "not_found", Message: err.Error()})
	case errors.Is(err, domain.ErrActiveFastExists),
		errors.Is(err, domain.ErrActiveSessionExists),
		errors.Is(err, domain.ErrActiveStreakExists),
		errors.Is(err, domain.ErrAlreadyLoggedToday):
		JSON(w, http.StatusConflict, ErrorResponse{Error: "conflict", Message: err.Error()})
	case errors.Is(err, domain.ErrNoActiveFast),
		errors.Is(err, domain.ErrNoActiveSession),
		errors.Is(err, domain.ErrNoActiveStreak):
		JSON(w, http.StatusBadRequest, ErrorResponse{Error: "bad_request", Message: err.Error()})
	case errors.Is(err, domain.ErrInvalidWorkoutType),
		errors.Is(err, domain.ErrInvalidDuration),
		errors.Is(err, domain.ErrInvalidMoodValue),
		errors.Is(err, domain.ErrInvalidEnergyLevel):
		JSON(w, http.StatusUnprocessableEntity, ErrorResponse{Error: "validation_error", Message: err.Error()})
	default:
		JSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal_error", Message: "An unexpected error occurred"})
	}
}

// DecodeJSON decodes a JSON request body into the given target.
func DecodeJSON(r *http.Request, target interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}
