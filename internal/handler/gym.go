package handler

import (
	"log"
	"net/http"

	"github.com/akshay/productiv-backend/internal/service"
)

// GymHandler handles gym HTTP requests.
type GymHandler struct {
	svc *service.GymService
}

// NewGymHandler creates a new GymHandler.
func NewGymHandler(svc *service.GymService) *GymHandler {
	return &GymHandler{svc: svc}
}

// GetStats returns gym statistics.
func (h *GymHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetStats(r.Context(), defaultUserID)
	if err != nil {
		log.Printf("error getting gym stats: %v", err)
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, stats)
}

// LogWorkout records a gym workout.
func (h *GymHandler) LogWorkout(w http.ResponseWriter, r *http.Request) {
	var req service.LogWorkoutRequest
	if err := DecodeJSON(r, &req); err != nil {
		JSON(w, http.StatusBadRequest, ErrorResponse{Error: "bad_request", Message: "Invalid request body"})
		return
	}

	session, err := h.svc.LogWorkout(r.Context(), defaultUserID, req)
	if err != nil {
		log.Printf("error logging workout: %v", err)
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusCreated, session)
}
