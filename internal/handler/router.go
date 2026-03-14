package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// NewRouter creates and configures the API router with all routes.
func NewRouter(
	dashboard *DashboardHandler,
	fasting *FastingHandler,
	gym *GymHandler,
	meditation *MeditationHandler,
	retention *RetentionHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		// Dashboard
		r.Get("/dashboard", dashboard.GetDashboard)

		// Fasting
		r.Route("/fasting", func(r chi.Router) {
			r.Get("/stats", fasting.GetStats)
			r.Get("/protocols", fasting.GetProtocols)
			r.Post("/start", fasting.StartFast)
			r.Post("/end", fasting.EndFast)
		})

		// Gym
		r.Route("/gym", func(r chi.Router) {
			r.Get("/stats", gym.GetStats)
			r.Post("/log", gym.LogWorkout)
		})

		// Meditation
		r.Route("/meditation", func(r chi.Router) {
			r.Get("/stats", meditation.GetStats)
			r.Post("/start", meditation.StartSession)
			r.Post("/end", meditation.EndSession)
		})

		// Retention
		r.Route("/retention", func(r chi.Router) {
			r.Get("/stats", retention.GetStats)
			r.Post("/start", retention.StartTracking)
			r.Post("/reset", retention.ResetCounter)
		})
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return r
}
