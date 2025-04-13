package main

import (
	middlewares "app05/internal/api/middleware"
	"app05/internal/api/routes"
	"app05/internal/core/application/contracts"
	"app05/internal/infrastructure/cache"
	"app05/internal/infrastructure/config"
	"app05/internal/infrastructure/logger"
	"app05/internal/infrastructure/rate_limiter"
	"app05/internal/infrastructure/scheduler"
	"app05/internal/infrastructure/scheduler/jobs"
	"app05/internal/infrastructure/services"
	"app05/internal/infrastructure/storage"
	"app05/pkg/appErrors"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Panicln("\033[31mError: Unable to load .env file, application will stop\033[0m")
	}

	// LOAD CONFIGURATION
	cfg := config.LoadConfig()

	//INITIALIZE LOGGER
	myLogger := logger.NewZapLogger()
	defer func(logger contracts.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Fatal("Failed to close logger: %v", err)
		}
	}(myLogger)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// DATABASE CONNECTION
	dbConn, err := config.InitDB(cfg)
	if err != nil {
		log.Panic(err)
	}
	defer func(dbConn *sql.DB) {
		err := dbConn.Close()
		if err != nil {
			myLogger.Fatal("Failed to close database connection: %v", err)
		}
	}(dbConn)

	// INITIALIZE REDIS CACHE
	myLogger.Info("Connecting to Redis...")
	redisCache, err := cache.NewSessionCache(cfg.Redis.URL, myLogger)
	if err != nil {
		myLogger.Fatal("Failed to initialize Redis", "error", err)
	}
	myLogger.Info("Successfully connected to Redis",
		"url", cfg.Redis.URL,
		"db", cfg.Redis.DB)

	// STORAGE INITIALIZATION
	store := storage.NewStorage(dbConn)

	// REGISTER SERVICES
	authService := services.NewAuthService(store.User, store.Session, redisCache, myLogger)
	userService := services.NewUserService(store.User, myLogger)
	serverService := services.NewServerStatusService(cfg.AppVersion, cfg.Env)
	postService := services.NewPostService(store.Post)

	//RATE LIMITER
	rL := rate_limiter.NewFixedWindowRateLimiter(cfg.RateLimiter.RequestPerTimeFrame, cfg.RateLimiter.TimeFrame)

	// INITIALIZE SCHEDULER
	newScheduler := scheduler.NewScheduler(myLogger)

	// Add session cleanup job
	cleanupJob := jobs.NewSessionCleanupJob(
		store.Session,
		myLogger,
		24*time.Hour, // Run once per day
	)
	newScheduler.AddJob(cleanupJob)

	// Create context for graceful shutdown
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Start newScheduler
	fmt.Println("Starting scheduler...")
	newScheduler.Start(ctx)

	//INITIALIZE THE ROUTER AND REGISTER MIDDLEWARE STACK
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	// CORS middlewares
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Use this to allow specific origin hosts
		//AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Use(middlewares.RateLimiterMiddleware(rL, cfg.RateLimiter.Enabled))

	fileServer := http.FileServer(http.Dir("uploads"))
	router.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServer))

	// ROUTES
	router.Route("/api/v1", func(r chi.Router) {
		routes.RegisterServerStatusRoutes(r, serverService, myLogger)
		routes.RegisterAuthRoutes(r, redisCache, authService, myLogger)
		routes.RegisterUserRoutes(r, redisCache, userService, myLogger)
		routes.RegisterPostRoutes(r, postService, myLogger)

	})

	// Custom 404 handler
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		Error := appErrors.New(appErrors.CodeNotFound, "Resource not found")
		appErrors.HandleError(w, Error, myLogger)
	})

	// Server
	server := &http.Server{
		Addr:         cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		myLogger.Info("Starting server on port " + cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			myLogger.Fatal("Server stopped"+err.Error(), err)
		}
	}()

	myLogger.Info("Server started")
	<-quit // Block until signal is received

	myLogger.Info("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		myLogger.Error("Server shutdown failed", err)
	}

	myLogger.Info("Server gracefully stopped")
}

func runMigrations(logger contracts.Logger) {
	cmd := exec.Command("make", "migrate-up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Info("Running database migrations...")
	if err := cmd.Run(); err != nil {
		logger.Fatal("Failed to run migrations: %v", err)
	}
	logger.Info("Database migrations complete.")
}
