package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"go-api-tech-challenge/internal/config"
	"go-api-tech-challenge/internal/database"
	"go-api-tech-challenge/internal/routes"
	"go-api-tech-challenge/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("Startup failed. err: %v", err)
	}
}

func run(ctx context.Context) error {

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("[in run]: %w", err)
	}

	logger := httplog.NewLogger("user-microservice", httplog.Options{
		LogLevel:        cfg.LogLevel,
		JSON:            false,
		Concise:         true,
		ResponseHeaders: false,
	})


	db, err := database.New(
		ctx,
		fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.DBHost,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			cfg.DBPort,
		),
		logger,
		time.Duration(cfg.DBRetryDuration)*time.Second,
	)
	if err != nil {
		return fmt.Errorf("[in run]: %w", err)
	}
  
	defer func() {
		if err = db.Close(); err != nil {
			logger.Error("Error closing db connection", "err", err)
		}
	}()

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:	[]string{"*"},
		AllowedMethods:	[]string{"GET". "POST", "PUT", "DELETE"},
		MaxAge: 300,
	}))


	svs := service.NewService(db)
	routes.RegisterRoutes(router, logger, svs, routes.WithRegisterHealthRoute(true))

	if cfg.HTTPUseSwagger {
		swagger.RunSwaggerUI(r, logger, cfg.HTTPDomain+cfg.HTTPPort)
	}

	srv := &http.Server{
		Addr:    cfg.HTTPDomain + cfg.HTTPPort,
		IdleTimeout: time.Minute,
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout: 500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
		Handler: r,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		<-sig

		fmt.Println()
		logger.Info("Shutdown signal received. Shutting down server...")

		shutdownCtx, err := context.WithTimeout(
			serverCtx, time.Duration(cfg.HTTPShutdownTimeout)*time.Second,
		)
		if err != nil {
			log.Fatalf("Error creating context.WithTimeout. err: %v", err)
		}

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				logger.Error("Shutdown timeout exceeded. Forcing shutdown...")
			}
		}()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Error shutting down server. err: %v", err)
		}
		serverStopCtx()
	}()

	logger.Info(fmt.Sprintf("Server started at %s", serverInstance.Addr))
	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err 
	}

	<-serverCtx.Done()
	logger.Info("Shutdown complete. Exiting...")
	return nil
}

