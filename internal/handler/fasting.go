package handler

import (
	"log"
	"net/http"

	"github.com/akshay/productiv-backend/internal/service"
)

// FastingHandler handles fasting HTTP requests.
type FastingHandler struct {
	svc service.FastingServiceInterface
}

// NewFastingHandler creates a new FastingHandler.
func NewFastingHandler(svc service.FastingServiceInterface) *FastingHandler {
	return &FastingHandler{svc: svc}
}

// GetStats returns fasting statistics and active session.
func (h *FastingHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetStats(r.Context(), defaultUserID)
	if err != nil {
		log.Printf("error getting fasting stats: %v", err)
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, stats)
}

// StartFast begins a new fasting session.
func (h *FastingHandler) StartFast(w http.ResponseWriter, r *http.Request) {
	var req service.StartFastRequest
	if err := DecodeJSON(r, &req); err != nil {
		JSON(w, http.StatusBadRequest, ErrorResponse{Error: "bad_request", Message: "Invalid request body"})
		return
	}

	session, err := h.svc.StartFast(r.Context(), defaultUserID, req)
	if err != nil {
		log.Printf("error starting fast: %v", err)
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusCreated, session)
}

// EndFast completes the active fasting session.
func (h *FastingHandler) EndFast(w http.ResponseWriter, r *http.Request) {
	session, err := h.svc.EndFast(r.Context(), defaultUserID)
	if err != nil {
		log.Printf("error ending fast: %v", err)
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, session)
}

// GetProtocols returns available fasting protocols.
func (h *FastingHandler) GetProtocols(w http.ResponseWriter, r *http.Request) {
	protocols := h.svc.GetProtocols()
	JSON(w, http.StatusOK, protocols)
}
