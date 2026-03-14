package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akshay/productiv-backend/internal/domain"
	svcmocks "github.com/akshay/productiv-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDashboardHandler_GetDashboard_Success(t *testing.T) {
	mockSvc := new(svcmocks.MockDashboardService)
	expected := &domain.DashboardData{
		DisciplineScore: 10,
		Pillars: []domain.PillarSummary{
			{Type: domain.PillarFasting, Name: "Fasting", CurrentStreak: 5, HasActivityToday: true},
			{Type: domain.PillarGym, Name: "Gym", CurrentStreak: 5, HasActivityToday: true},
		},
		TodaysFocus: "All systems active. Stay the course.",
	}
	mockSvc.On("GetDashboard", mock.Anything, int64(1)).Return(expected, nil)

	h := NewDashboardHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard", nil)
	w := httptest.NewRecorder()
	h.GetDashboard(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got domain.DashboardData
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, 10, got.DisciplineScore)
	assert.Len(t, got.Pillars, 2)
	assert.Equal(t, "All systems active. Stay the course.", got.TodaysFocus)
	mockSvc.AssertExpectations(t)
}

func TestDashboardHandler_GetDashboard_ServiceError(t *testing.T) {
	mockSvc := new(svcmocks.MockDashboardService)
	mockSvc.On("GetDashboard", mock.Anything, int64(1)).Return(nil, errors.New("unexpected error"))

	h := NewDashboardHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard", nil)
	w := httptest.NewRecorder()
	h.GetDashboard(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var got ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "internal_error", got.Error)
}
