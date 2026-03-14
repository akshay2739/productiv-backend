package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/akshay/productiv-backend/internal/service"
	svcmocks "github.com/akshay/productiv-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMeditationHandler_GetStats_Success(t *testing.T) {
	mockSvc := new(svcmocks.MockMeditationService)
	expected := &domain.MeditationStats{
		CurrentStreak:  7,
		TotalMinutes:   200.5,
		TotalSessions:  15,
		AverageSession: 13.37,
	}
	mockSvc.On("GetStats", mock.Anything, int64(1)).Return(expected, nil)

	h := NewMeditationHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/meditation/stats", nil)
	w := httptest.NewRecorder()
	h.GetStats(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got domain.MeditationStats
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, 7, got.CurrentStreak)
	assert.Equal(t, 200.5, got.TotalMinutes)
	assert.Equal(t, 15, got.TotalSessions)
	mockSvc.AssertExpectations(t)
}

func TestMeditationHandler_GetStats_ServiceError(t *testing.T) {
	mockSvc := new(svcmocks.MockMeditationService)
	mockSvc.On("GetStats", mock.Anything, int64(1)).Return(nil, domain.ErrNotFound)

	h := NewMeditationHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/meditation/stats", nil)
	w := httptest.NewRecorder()
	h.GetStats(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMeditationHandler_StartSession_Success(t *testing.T) {
	mockSvc := new(svcmocks.MockMeditationService)
	mood := 3
	expected := &domain.MeditationSession{
		ID:            1,
		UserID:        1,
		TargetMinutes: 10,
		StartTime:     time.Now().UTC(),
		MoodBefore:    &mood,
	}
	mockSvc.On("StartSession", mock.Anything, int64(1), service.StartSessionRequest{
		TargetMinutes: 10,
		MoodBefore:    &mood,
	}).Return(expected, nil)

	h := NewMeditationHandler(mockSvc)
	body, _ := json.Marshal(map[string]interface{}{"target_minutes": 10, "mood_before": 3})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/meditation/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.StartSession(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var got domain.MeditationSession
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, 10, got.TargetMinutes)
	mockSvc.AssertExpectations(t)
}

func TestMeditationHandler_StartSession_BadJSON(t *testing.T) {
	h := NewMeditationHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/meditation/start", bytes.NewReader([]byte("bad")))
	w := httptest.NewRecorder()
	h.StartSession(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMeditationHandler_StartSession_Conflict(t *testing.T) {
	mockSvc := new(svcmocks.MockMeditationService)
	mockSvc.On("StartSession", mock.Anything, int64(1), mock.Anything).Return(nil, domain.ErrActiveSessionExists)

	h := NewMeditationHandler(mockSvc)
	body, _ := json.Marshal(map[string]interface{}{"target_minutes": 10})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/meditation/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.StartSession(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestMeditationHandler_EndSession_Success(t *testing.T) {
	mockSvc := new(svcmocks.MockMeditationService)
	dur := 10.5
	now := time.Now().UTC()
	mood := 5
	expected := &domain.MeditationSession{
		ID:             1,
		EndTime:        &now,
		ActualDuration: &dur,
		MoodAfter:      &mood,
	}
	mockSvc.On("EndSession", mock.Anything, int64(1), service.EndSessionRequest{MoodAfter: &mood}).Return(expected, nil)

	h := NewMeditationHandler(mockSvc)
	body, _ := json.Marshal(map[string]interface{}{"mood_after": 5})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/meditation/end", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.EndSession(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestMeditationHandler_EndSession_EmptyBody(t *testing.T) {
	mockSvc := new(svcmocks.MockMeditationService)
	dur := 5.0
	now := time.Now().UTC()
	expected := &domain.MeditationSession{
		ID:             1,
		EndTime:        &now,
		ActualDuration: &dur,
	}
	mockSvc.On("EndSession", mock.Anything, int64(1), service.EndSessionRequest{}).Return(expected, nil)

	h := NewMeditationHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/meditation/end", nil)
	w := httptest.NewRecorder()
	h.EndSession(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestMeditationHandler_EndSession_NoActiveSession(t *testing.T) {
	mockSvc := new(svcmocks.MockMeditationService)
	mockSvc.On("EndSession", mock.Anything, int64(1), mock.Anything).Return(nil, domain.ErrNoActiveSession)

	h := NewMeditationHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/meditation/end", nil)
	w := httptest.NewRecorder()
	h.EndSession(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
