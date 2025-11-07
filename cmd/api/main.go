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

	"github.com/abdullinmm/user-management-api/internal/config"
	"github.com/abdullinmm/user-management-api/internal/middleware"
	jwtpkg "github.com/abdullinmm/user-management-api/internal/pkg/jwt"
	"github.com/abdullinmm/user-management-api/internal/repository/postgresql"
	"github.com/abdullinmm/user-management-api/internal/usecase"

	"github.com/go-chi/chi/v5"
	middleware2 "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	dbPool, err := initDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	log.Println("Successfully connected to database")

	// Initialize repositories
	userRepo := postgresql.NewUserRepository(dbPool)
	taskRepo := postgresql.NewTaskRepository(dbPool)
	balanceRepo := postgresql.NewBalanceRepository(dbPool)
	transactionRepo := postgresql.NewTransactionRepository(dbPool)
	userTaskRepo := postgresql.NewUserTaskRepository(dbPool)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo, balanceRepo, transactionRepo, cfg.ReferralBonus, cfg.RefereeBonus)
	taskUseCase := usecase.NewTaskUseCase(taskRepo, userTaskRepo, balanceRepo, transactionRepo)
	balanceUseCase := usecase.NewBalanceUseCase(balanceRepo, transactionRepo)

	// Initialize JWT manager
	jwtManager := jwtpkg.NewManager(cfg.JWTSecret)

	// Initialize HTTP router
	router := setupRouter(userUseCase, taskUseCase, balanceUseCase, jwtManager)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on port %s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

// initDatabase initializes database connection pool
func initDatabase(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Configure pool settings
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}

// setupRouter configures HTTP routes and middleware
func setupRouter(
	userUC *usecase.UserUseCase,
	taskUC *usecase.TaskUseCase,
	balanceUC *usecase.BalanceUseCase,
	jwtManager *jwtpkg.Manager,
) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware2.RequestID)
	r.Use(middleware2.RealIP)
	r.Use(middleware2.Logger)
	r.Use(middleware2.Recoverer)
	r.Use(middleware2.Timeout(60 * time.Second))

	// Health check (no auth required)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Public routes (no auth required)
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"registration endpoint - will be implemented in handlers"}`))
		})

		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"login endpoint - will be implemented in handlers"}`))
		})
	})

	// Protected routes (JWT auth required)
	r.Route("/api/v1/users", func(r chi.Router) {
		// Apply JWT middleware to all routes in this group
		r.Use(middleware.Auth(jwtManager))

		// GET /users/{id}/status - get user status
		r.Get("/{id}/status", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"user status endpoint - protected by JWT"}`))
		})

		// GET /users/leaderboard - get top users
		r.Get("/leaderboard", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"leaderboard endpoint - protected by JWT"}`))
		})

		// POST /users/{id}/task/complete - complete task
		r.Post("/{id}/task/complete", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"complete task endpoint - protected by JWT"}`))
		})

		// POST /users/{id}/referrer - set referrer
		r.Post("/{id}/referrer", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"set referrer endpoint - protected by JWT"}`))
		})
	})

	return r
}
