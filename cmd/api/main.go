package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mifinca-backend/internal/config"
	animalApp "mifinca-backend/internal/application/animal"
	authApp "mifinca-backend/internal/application/auth"
	healthApp "mifinca-backend/internal/application/health"
	movementApp "mifinca-backend/internal/application/movement"
	httpDelivery "mifinca-backend/internal/delivery/http"
	"mifinca-backend/internal/infrastructure/auth"
	"mifinca-backend/internal/infrastructure/database"
	"mifinca-backend/internal/infrastructure/repository"
	"mifinca-backend/internal/infrastructure/websocket"
)

func main() {
	// ── Load Configuration ────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: failed to load config: %v", err)
	}

	// ── Connect to Database ───────────────────────────────────────────────────
	ctx := context.Background()
	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("FATAL: failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("✅ Connected to PostgreSQL")

	// ── Infrastructure ────────────────────────────────────────────────────────
	jwtMaker := auth.NewJWTMaker(cfg.JWTSecret, cfg.JWTAccessDuration, cfg.JWTRefreshDuration)

	// ── Repositories (DI Manual) ──────────────────────────────────────────────
	userRepo := repository.NewUserRepository(pool)
	tenantRepo := repository.NewTenantRepository(pool)
	deviceRepo := repository.NewDeviceRepository(pool)
	farmRepo := repository.NewFarmRepository(pool)
	animalRepo := repository.NewAnimalRepository(pool)
	movementRepo := repository.NewMovementRepository(pool)
	healthRepo := repository.NewHealthRepository(pool)
	
	hub := websocket.NewMockHub()

	// ── Application Services ──────────────────────────────────────────────────
	authService := authApp.NewService(userRepo, tenantRepo, deviceRepo, farmRepo, jwtMaker)
	movementService := movementApp.NewService(movementRepo, hub)
	healthService := healthApp.NewService(healthRepo, hub)
	animalService := animalApp.NewService(animalRepo, farmRepo, movementRepo, hub)

	// ── HTTP Router ───────────────────────────────────────────────────────────
	router := httpDelivery.NewRouter(pool, jwtMaker, authService, animalService, healthService, movementService)

	// ── Server ────────────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ── Graceful Shutdown ─────────────────────────────────────────────────────
	go func() {
		log.Printf("🚀 MiFinca API starting on port %s (env=%s)", cfg.Port, cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("FATAL: server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⏳ Shutting down gracefully...")
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("FATAL: server forced to shutdown: %v", err)
	}
	log.Println("✅ Server stopped")
}
