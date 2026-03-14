package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/akshay/productiv-backend/internal/handler"
	"github.com/akshay/productiv-backend/internal/service"
	svcmocks "github.com/akshay/productiv-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGymHandler_GetStats_Success(t *testing.T) {
	mockSvc := new(svcmocks.MockGymService)
	expected := &domain.GymStats{
		CurrentStreak: 5,
		LongestStreak: 10,
		WorkoutsWeek:  3,
		WorkoutsMonth: 12,
		WeeklyGoal:    5,
		LoggedToday:   true,
	}
	mockSvc.On("GetStats", mock.Anything, int64(1)).Return(expected, nil)

	h := handler.NewGymHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/gym/stats", nil)
	w := httptest.NewRecorder()
	h.GetStats(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got domain.GymStats
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, 5, got.CurrentStreak)
	assert.Equal(t, 3, got.WorkoutsWeek)
	assert.True(t, got.LoggedToday)
	mockSvc.AssertExpectations(t)
}

func TestGymHandler_GetStats_ServiceError(t *testing.T) {
	mockSvc := new(svcmocks.MockGymService)
	mockSvc.On("GetStats", mock.Anything, int64(1)).Return(nil, domain.ErrNotFound)

	h := handler.NewGymHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/gym/stats", nil)
	w := httptest.NewRecorder()
	h.GetStats(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGymHandler_LogWorkout_Success(t *testing.T) {
	mockSvc := new(svcmocks.MockGymService)
	expected := &domain.GymSession{
		ID:          1,
		UserID:      1,
		WorkoutType: domain.WorkoutStrength,
	}
	mockSvc.On("LogWorkout", mock.Anything, int64(1), service.LogWorkoutRequest{WorkoutType: "strength"}).Return(expected, nil)

	h := handler.NewGymHandler(mockSvc)
	body, _ := json.Marshal(map[string]string{"workout_type": "strength"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gym/log", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.LogWorkout(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var got domain.GymSession
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, domain.WorkoutStrength, got.WorkoutType)
	mockSvc.AssertExpectations(t)
}

func TestGymHandler_LogWorkout_BadJSON(t *testing.T) {
	mockSvc := new(svcmocks.MockGymService)
	h := handler.NewGymHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gym/log", bytes.NewReader([]byte("{invalid")))
	w := httptest.NewRecorder()
	h.LogWorkout(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGymHandler_LogWorkout_Conflict(t *testing.T) {
	mockSvc := new(svcmocks.MockGymService)
	mockSvc.On("LogWorkout", mock.Anything, int64(1), mock.Anything).Return(nil, domain.ErrAlreadyLoggedToday)

	h := handler.NewGymHandler(mockSvc)
	body, _ := json.Marshal(map[string]string{"workout_type": "cardio"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gym/log", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.LogWorkout(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var got handler.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "conflict", got.Error)
}

func TestGymHandler_LogWorkout_ValidationError(t *testing.T) {
	mockSvc := new(svcmocks.MockGymService)
	mockSvc.On("LogWorkout", mock.Anything, int64(1), mock.Anything).Return(nil, domain.ErrInvalidWorkoutType)

	h := handler.NewGymHandler(mockSvc)
	body, _ := json.Marshal(map[string]string{"workout_type": "boxing"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gym/log", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.LogWorkout(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}
