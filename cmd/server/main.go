package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akshay/productiv-backend/config"
	"github.com/akshay/productiv-backend/internal/handler"
	"github.com/akshay/productiv-backend/internal/middleware"
	"github.com/akshay/productiv-backend/internal/repository/postgres"
	"github.com/akshay/productiv-backend/internal/service"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := postgres.NewPool(ctx, cfg.DBConnString())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("connected to database")

	// Initialize repositories
	userRepo := postgres.NewUserRepo(pool)
	pillarRepo := postgres.NewPillarRepo(pool)
	fastingRepo := postgres.NewFastingRepo(pool)
	gymRepo := postgres.NewGymRepo(pool)
	meditationRepo := postgres.NewMeditationRepo(pool)
	retentionRepo := postgres.NewRetentionRepo(pool)

	// Initialize services
	fastingSvc := service.NewFastingService(fastingRepo, userRepo)
	gymSvc := service.NewGymService(gymRepo, userRepo)
	meditationSvc := service.NewMeditationService(meditationRepo, userRepo)
	retentionSvc := service.NewRetentionService(retentionRepo, userRepo)
	dashboardSvc := service.NewDashboardService(pillarRepo, fastingRepo, gymRepo, meditationRepo, retentionRepo, userRepo)

	// Initialize handlers
	dashboardHandler := handler.NewDashboardHandler(dashboardSvc)
	fastingHandler := handler.NewFastingHandler(fastingSvc)
	gymHandler := handler.NewGymHandler(gymSvc)
	meditationHandler := handler.NewMeditationHandler(meditationSvc)
	retentionHandler := handler.NewRetentionHandler(retentionSvc)

	// Setup router
	router := handler.NewRouter(dashboardHandler, fastingHandler, gymHandler, meditationHandler, retentionHandler)

	// Apply global middleware
	router.Use(middleware.CORS(cfg.CORSOrigin))
	router.Use(middleware.Logging)

	// Start server
	srv := &http.Server{
		Addr:         cfg.ServerAddr(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("server shutdown failed: %v", err)
		}
	}()

	log.Printf("server starting on %s", cfg.ServerAddr())
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
	log.Println("server stopped")
}
