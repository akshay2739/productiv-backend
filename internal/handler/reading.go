package handler

import (
	"log"
	"net/http"

	"github.com/akshay/productiv-backend/internal/service"
)

// ReadingHandler handles reading HTTP requests.
type ReadingHandler struct {
	svc service.ReadingServiceInterface
}

// NewReadingHandler creates a new ReadingHandler.
func NewReadingHandler(svc service.ReadingServiceInterface) *ReadingHandler {
	return &ReadingHandler{svc: svc}
}

// GetStats returns reading statistics.
func (h *ReadingHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetStats(r.Context(), defaultUserID)
	if err != nil {
		log.Printf("error getting reading stats: %v", err)
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, stats)
}

// LogReading records a reading session.
func (h *ReadingHandler) LogReading(w http.ResponseWriter, r *http.Request) {
	var req service.LogReadingRequest
	if err := DecodeJSON(r, &req); err != nil {
		JSON(w, http.StatusBadRequest, ErrorResponse{Error: "bad_request", Message: "Invalid request body"})
		return
	}

	session, err := h.svc.LogReading(r.Context(), defaultUserID, req)
	if err != nil {
		log.Printf("error logging reading: %v", err)
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusCreated, session)
}
