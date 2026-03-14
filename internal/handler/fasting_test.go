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

func TestFastingHandler_GetStats_Success(t *testing.T) {
	mockSvc := new(svcmocks.MockFastingService)
	expected := &domain.FastingStats{
		CurrentStreak:   3,
		LongestStreak:   5,
		AverageDuration: 16.5,
		TotalFasts:      10,
		CalendarDays:    []domain.FastingCalendarDay{{Date: "2026-03-14", HasActivity: true, IsToday: true}},
	}
	mockSvc.On("GetStats", mock.Anything, int64(1)).Return(expected, nil)

	h := handler.NewFastingHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/fasting/stats", nil)
	w := httptest.NewRecorder()
	h.GetStats(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got domain.FastingStats
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, 3, got.CurrentStreak)
	assert.Equal(t, 10, got.TotalFasts)
	assert.Equal(t, 16.5, got.AverageDuration)
	mockSvc.AssertExpectations(t)
}

func TestFastingHandler_GetStats_ServiceError(t *testing.T) {
	mockSvc := new(svcmocks.MockFastingService)
	mockSvc.On("GetStats", mock.Anything, int64(1)).Return(nil, domain.ErrNotFound)

	h := handler.NewFastingHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/fasting/stats", nil)
	w := httptest.NewRecorder()
	h.GetStats(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestFastingHandler_StartFast_Success(t *testing.T) {
	mockSvc := new(svcmocks.MockFastingService)
	expected := &domain.FastingSession{
		ID:          1,
		UserID:      1,
		Protocol:    "16:8",
		TargetHours: 16,
		StartTime:   time.Now().UTC(),
	}
	mockSvc.On("StartFast", mock.Anything, int64(1), service.StartFastRequest{Protocol: "16:8"}).Return(expected, nil)

	h := handler.NewFastingHandler(mockSvc)
	body, _ := json.Marshal(map[string]string{"protocol": "16:8"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/fasting/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.StartFast(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var got domain.FastingSession
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "16:8", got.Protocol)
	mockSvc.AssertExpectations(t)
}

func TestFastingHandler_StartFast_BadJSON(t *testing.T) {
	mockSvc := new(svcmocks.MockFastingService)
	h := handler.NewFastingHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/fasting/start", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()
	h.StartFast(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFastingHandler_StartFast_Conflict(t *testing.T) {
	mockSvc := new(svcmocks.MockFastingService)
	mockSvc.On("StartFast", mock.Anything, int64(1), mock.Anything).Return(nil, domain.ErrActiveFastExists)

	h := handler.NewFastingHandler(mockSvc)
	body, _ := json.Marshal(map[string]string{"protocol": "16:8"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/fasting/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.StartFast(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var got handler.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "conflict", got.Error)
}

func TestFastingHandler_StartFast_ValidationError(t *testing.T) {
	mockSvc := new(svcmocks.MockFastingService)
	mockSvc.On("StartFast", mock.Anything, int64(1), mock.Anything).Return(nil, domain.ErrInvalidDuration)

	h := handler.NewFastingHandler(mockSvc)
	body, _ := json.Marshal(map[string]string{"protocol": "bad"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/fasting/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.StartFast(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestFastingHandler_EndFast_Success(t *testing.T) {
	mockSvc := new(svcmocks.MockFastingService)
	dur := 16.5
	now := time.Now().UTC()
	expected := &domain.FastingSession{
		ID:             1,
		Protocol:       "16:8",
		EndTime:        &now,
		ActualDuration: &dur,
		TargetReached:  true,
	}
	mockSvc.On("EndFast", mock.Anything, int64(1)).Return(expected, nil)

	h := handler.NewFastingHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/fasting/end", nil)
	w := httptest.NewRecorder()
	h.EndFast(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestFastingHandler_EndFast_NoActiveFast(t *testing.T) {
	mockSvc := new(svcmocks.MockFastingService)
	mockSvc.On("EndFast", mock.Anything, int64(1)).Return(nil, domain.ErrNoActiveFast)

	h := handler.NewFastingHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/fasting/end", nil)
	w := httptest.NewRecorder()
	h.EndFast(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFastingHandler_GetProtocols_Success(t *testing.T) {
	mockSvc := new(svcmocks.MockFastingService)
	protocols := domain.AvailableFastingProtocols()
	mockSvc.On("GetProtocols").Return(protocols)

	h := handler.NewFastingHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/fasting/protocols", nil)
	w := httptest.NewRecorder()
	h.GetProtocols(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got []domain.FastingProtocol
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Len(t, got, 4)
	mockSvc.AssertExpectations(t)
}
