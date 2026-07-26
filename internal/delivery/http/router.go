package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	animalApp "mifinca-backend/internal/application/animal"
	authApp "mifinca-backend/internal/application/auth"
	healthApp "mifinca-backend/internal/application/health"
	movementApp "mifinca-backend/internal/application/movement"
	"mifinca-backend/internal/delivery/http/handler"
	"mifinca-backend/internal/delivery/http/httputil"
	"mifinca-backend/internal/delivery/http/middleware"
	infraAuth "mifinca-backend/internal/infrastructure/auth"
)

// NewRouter creates the Chi router with all middleware and route groups.
func NewRouter(
	pool *pgxpool.Pool,
	jwtMaker *infraAuth.JWTMaker,
	authService *authApp.Service,
	animalService *animalApp.Service,
	healthService *healthApp.Service,
	movementService *movementApp.Service,
	userService *userApp.UserService,
) http.Handler {
	r := chi.NewRouter()

	// ── Global Middleware ──────────────────────────────────────────────────────
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Compress(5))
	r.Use(middleware.CORS())

	// ── Health Check ──────────────────────────────────────────────────────────
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// ── Handlers ──────────────────────────────────────────────────────────────
	authHandler := handler.NewAuthHandler(authService)
	animalHandler := handler.NewAnimalHandler(animalService)
	healthHandler := handler.NewHealthHandler(healthService)
	movementHandler := handler.NewMovementHandler(movementService)
	userHandler := handler.NewUserHandler(userService)

	// ── API v1 Routes ─────────────────────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {

		// ── Public routes (no JWT required) ───────────────────────────────
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", authHandler.Register)
			r.Post("/auth/login", authHandler.Login)
			r.Post("/auth/refresh", authHandler.Refresh)
		})

		// ── Protected routes (JWT required) ───────────────────────────────
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtMaker))
			r.Use(middleware.Tenant(pool))

			// Users
			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.List)
				r.Post("/", userHandler.Create)
				r.Get("/{id}", userHandler.GetByID)
				r.Put("/{id}", userHandler.Update)
				r.Delete("/{id}", userHandler.Delete)
			})

			// Animals
			r.Route("/farms/{farm_id}/animals", func(r chi.Router) {
				r.Get("/", animalHandler.ListAnimals)
				r.Post("/", animalHandler.Create)
			})
			r.Route("/animals/{id}", func(r chi.Router) {
				r.Get("/", animalHandler.GetAnimal)
				r.Put("/", animalHandler.Update)
				r.Patch("/status", animalHandler.ChangeStatus)
				r.Delete("/", animalHandler.SoftDelete)
				r.Get("/pedigree", animalHandler.GetPedigree)
				r.Get("/offspring", notImplemented)  // GetOffspring — UC-05.2
			})

			// Health Events
			r.Get("/animals/{id}/health", notImplemented)          // ListHealth — UC-03.3
			r.Post("/animals/{id}/health", healthHandler.Create)   // CreateHealth — UC-03.1
			r.Post("/farms/{farm_id}/health/bulk", notImplemented) // BulkHealth — UC-03.2

			// Movements
			r.Get("/animals/{id}/movements", notImplemented)        // ListMovements — UC-04.2
			r.Post("/animals/{id}/movements", movementHandler.Create) // CreateMovement — UC-04.1

			// Farms
			r.Get("/farms", notImplemented)  // ListFarms
			r.Post("/farms", notImplemented) // CreateFarm — UC-08.1

			// Reports
			r.Get("/reports/summary", notImplemented)     // HerdSummary — UC-07.1
			r.Get("/reports/vaccination", notImplemented)  // VaccinationReport — UC-07.2

			// Sync
			r.Post("/sync/push", notImplemented) // Push — UC-06.1
			r.Get("/sync/pull", notImplemented)   // Pull — UC-06.2

			// WebSocket
			r.Get("/ws", notImplemented) // ServeWS — Tiempo Real
		})
	})

	return r
}

// notImplemented is a placeholder handler for routes not yet implemented.
func notImplemented(w http.ResponseWriter, r *http.Request) {
	httputil.ErrorWithMessage(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "this endpoint is under construction")
}
