package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/akshay/productiv-backend/internal/handler"
	"github.com/akshay/productiv-backend/internal/service"
	svcmocks "github.com/akshay/productiv-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRetentionHandler_GetStats_Success(t *testing.T) {
	mockSvc := new(svcmocks.MockRetentionService)
	milestone := 30
	expected := &domain.RetentionStats{
		CurrentDayCount: 15,
		BestStreak:      20,
		NextMilestone:   &milestone,
		Milestones: []domain.Milestone{
			{Days: 7, Achieved: true},
			{Days: 14, Achieved: true},
			{Days: 30, Achieved: false},
		},
		ActiveStreak: &domain.RetentionStreak{ID: 1, StartDate: time.Now().UTC()},
	}
	mockSvc.On("GetStats", mock.Anything, int64(1)).Return(expected, nil)

	h := handler.NewRetentionHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/retention/stats", nil)
	w := httptest.NewRecorder()
	h.GetStats(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got domain.RetentionStats
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, 15, got.CurrentDayCount)
	assert.Equal(t, 20, got.BestStreak)
	assert.NotNil(t, got.NextMilestone)
	assert.Equal(t, 30, *got.NextMilestone)
	assert.Len(t, got.Milestones, 3)
	mockSvc.AssertExpectations(t)
}

func TestRetentionHandler_GetStats_ServiceError(t *testing.T) {
	mockSvc := new(svcmocks.MockRetentionService)
	mockSvc.On("GetStats", mock.Anything, int64(1)).Return(nil, domain.ErrNotFound)

	h := handler.NewRetentionHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/retention/stats", nil)
	w := httptest.NewRecorder()
	h.GetStats(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRetentionHandler_StartTracking_Success(t *testing.T) {
	mockSvc := new(svcmocks.MockRetentionService)
	expected := &domain.RetentionStreak{
		ID:        1,
		UserID:    1,
		StartDate: time.Now().UTC(),
	}
	mockSvc.On("StartTracking", mock.Anything, int64(1)).Return(expected, nil)

	h := handler.NewRetentionHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/retention/start", nil)
	w := httptest.NewRecorder()
	h.StartTracking(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var got domain.RetentionStreak
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, int64(1), got.ID)
	mockSvc.AssertExpectations(t)
}

func TestRetentionHandler_StartTracking_Conflict(t *testing.T) {
	mockSvc := new(svcmocks.MockRetentionService)
	mockSvc.On("StartTracking", mock.Anything, int64(1)).Return(nil, domain.ErrActiveStreakExists)

	h := handler.NewRetentionHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/retention/start", nil)
	w := httptest.NewRecorder()
	h.StartTracking(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRetentionHandler_ResetCounter_Success(t *testing.T) {
	mockSvc := new(svcmocks.MockRetentionService)
	reason := "testing"
	now := time.Now().UTC()
	expected := &domain.RetentionStreak{
		ID:        1,
		EndDate:   &now,
		DaysCount: 10,
		Reason:    &reason,
	}
	mockSvc.On("ResetCounter", mock.Anything, int64(1), service.ResetRequest{Reason: &reason}).Return(expected, nil)

	h := handler.NewRetentionHandler(mockSvc)
	body, _ := json.Marshal(map[string]string{"reason": "testing"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/retention/reset", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ResetCounter(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got domain.RetentionStreak
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, 10, got.DaysCount)
	assert.Equal(t, &reason, got.Reason)
	mockSvc.AssertExpectations(t)
}

func TestRetentionHandler_ResetCounter_EmptyBody(t *testing.T) {
	mockSvc := new(svcmocks.MockRetentionService)
	now := time.Now().UTC()
	expected := &domain.RetentionStreak{ID: 1, EndDate: &now, DaysCount: 5}
	mockSvc.On("ResetCounter", mock.Anything, int64(1), service.ResetRequest{}).Return(expected, nil)

	h := handler.NewRetentionHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/retention/reset", nil)
	w := httptest.NewRecorder()
	h.ResetCounter(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestRetentionHandler_ResetCounter_NoActiveStreak(t *testing.T) {
	mockSvc := new(svcmocks.MockRetentionService)
	mockSvc.On("ResetCounter", mock.Anything, int64(1), mock.Anything).Return(nil, domain.ErrNoActiveStreak)

	h := handler.NewRetentionHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/retention/reset", nil)
	w := httptest.NewRecorder()
	h.ResetCounter(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
