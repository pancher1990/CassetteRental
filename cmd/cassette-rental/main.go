package main

import (
	"CassetteRental/internal/config"
	addCassette "CassetteRental/internal/http-server/handlers/cassette/create"
	"CassetteRental/internal/http-server/handlers/cassette/setStatus"
	customerCreate "CassetteRental/internal/http-server/handlers/customer/create"
	"CassetteRental/internal/http-server/handlers/customer/setBalance"
	addFilm "CassetteRental/internal/http-server/handlers/film/create"
	findFilm "CassetteRental/internal/http-server/handlers/film/find"
	makeRent "CassetteRental/internal/http-server/handlers/rent/create"
	"CassetteRental/internal/storage/postgresql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	log := setupLogger(cfg.Env)
	log.Info("starting ", slog.String("env", cfg.Env))
	log.Info(cfg.Env)
	storage, err := postgresql.New(cfg.StorageDb)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Post("/customer/create", customerCreate.New(log, storage))
	router.Patch("/customer/balance/{customerId}", setBalance.New(log, storage))

	router.Post("/film/add", addFilm.New(log, storage))
	router.Get("/film/find", findFilm.New(log, storage))

	router.Post("/cassette/add", addCassette.New(log, storage))
	router.Patch("/cassette/available/{cassetteId}", setStatus.New(log, storage))

	router.Post("/rent/create", makeRent.New(log, storage))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")

	}

	log.Error("server stopped")
}
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		))
	case envDev:
		log = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		))
	case envProd:
		log = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelError},
		))
	}
	return log
}
