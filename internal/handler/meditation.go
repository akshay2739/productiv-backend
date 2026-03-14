package handler

import (
	"log"
	"net/http"

	"github.com/akshay/productiv-backend/internal/service"
)

// MeditationHandler handles meditation HTTP requests.
type MeditationHandler struct {
	svc *service.MeditationService
}

// NewMeditationHandler creates a new MeditationHandler.
func NewMeditationHandler(svc *service.MeditationService) *MeditationHandler {
	return &MeditationHandler{svc: svc}
}

// GetStats returns meditation statistics.
func (h *MeditationHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetStats(r.Context(), defaultUserID)
	if err != nil {
		log.Printf("error getting meditation stats: %v", err)
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, stats)
}

// StartSession begins a new meditation session.
func (h *MeditationHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	var req service.StartSessionRequest
	if err := DecodeJSON(r, &req); err != nil {
		JSON(w, http.StatusBadRequest, ErrorResponse{Error: "bad_request", Message: "Invalid request body"})
		return
	}

	session, err := h.svc.StartSession(r.Context(), defaultUserID, req)
	if err != nil {
		log.Printf("error starting meditation: %v", err)
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusCreated, session)
}

// EndSession completes the active meditation session.
func (h *MeditationHandler) EndSession(w http.ResponseWriter, r *http.Request) {
	var req service.EndSessionRequest
	if err := DecodeJSON(r, &req); err != nil {
		// Allow empty body for ending session
		req = service.EndSessionRequest{}
	}

	session, err := h.svc.EndSession(r.Context(), defaultUserID, req)
	if err != nil {
		log.Printf("error ending meditation: %v", err)
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, session)
}
