package handler

import (
	"log"
	"net/http"

	"github.com/akshay/productiv-backend/internal/service"
)

// RetentionHandler handles retention HTTP requests.
type RetentionHandler struct {
	svc service.RetentionServiceInterface
}

// NewRetentionHandler creates a new RetentionHandler.
func NewRetentionHandler(svc service.RetentionServiceInterface) *RetentionHandler {
	return &RetentionHandler{svc: svc}
}

// GetStats returns retention statistics.
func (h *RetentionHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetStats(r.Context(), defaultUserID)
	if err != nil {
		log.Printf("error getting retention stats: %v", err)
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, stats)
}

// StartTracking begins a new retention streak.
func (h *RetentionHandler) StartTracking(w http.ResponseWriter, r *http.Request) {
	streak, err := h.svc.StartTracking(r.Context(), defaultUserID)
	if err != nil {
		log.Printf("error starting retention: %v", err)
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusCreated, streak)
}

// ResetCounter ends the active retention streak.
func (h *RetentionHandler) ResetCounter(w http.ResponseWriter, r *http.Request) {
	var req service.ResetRequest
	if err := DecodeJSON(r, &req); err != nil {
		// Allow empty body for reset
		req = service.ResetRequest{}
	}

	streak, err := h.svc.ResetCounter(r.Context(), defaultUserID, req)
	if err != nil {
		log.Printf("error resetting retention: %v", err)
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, streak)
}
