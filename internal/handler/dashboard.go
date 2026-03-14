package handler

import (
	"log"
	"net/http"

	"github.com/akshay/productiv-backend/internal/service"
)

// DashboardHandler handles dashboard HTTP requests.
type DashboardHandler struct {
	svc *service.DashboardService
}

// NewDashboardHandler creates a new DashboardHandler.
func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// GetDashboard returns the dashboard data for the default user.
func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetDashboard(r.Context(), defaultUserID)
	if err != nil {
		log.Printf("error getting dashboard: %v", err)
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, data)
}
