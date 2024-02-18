package main

import (
	"CassetteRental/internal/config"
	"CassetteRental/internal/http-server/handlers/auth"
	addCassette "CassetteRental/internal/http-server/handlers/cassette/create"
	"CassetteRental/internal/http-server/handlers/cassette/setStatus"
	customerCreate "CassetteRental/internal/http-server/handlers/customer/create"
	"CassetteRental/internal/http-server/handlers/customer/setBalance"
	addFilm "CassetteRental/internal/http-server/handlers/film/create"
	findFilm "CassetteRental/internal/http-server/handlers/film/find"
	makeRent "CassetteRental/internal/http-server/handlers/rent/create"
	middlewareAuth "CassetteRental/internal/http-server/middleware/auth"
	token "CassetteRental/internal/lib/auth"
	"CassetteRental/internal/lib/hash"
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
		log.Error("failed to init storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	tokenManager, err := token.NewManager(cfg.Auth.SigningKey)
	if err != nil {
		log.Error("failed to init token manager", slog.String("error", err.Error()))
		os.Exit(1)
	}

	hasherForAuth, err := hash.NewHash(cfg.Auth.Salt)
	if err != nil {
		log.Error("failed to init hash maker", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Post("/auth", auth.New(log, storage, tokenManager, hasherForAuth))
	router.Post("/customer/create", customerCreate.New(log, storage, hasherForAuth))
	router.Get("/film/find", findFilm.New(log, storage))
	router.Group(func(r chi.Router) {
		r.Use(middlewareAuth.CheckToken(tokenManager))
		r.Post("/rent/create", makeRent.New(log, storage))
	})
	router.Route("/admin", func(r chi.Router) {
		r.Use(middleware.BasicAuth("cassette-rental",
			map[string]string{cfg.HTTPServer.AdminLogin: cfg.HTTPServer.AdminPassword}))
		r.Get("/film/find", findFilm.New(log, storage))

		r.Post("/cassette/add", addCassette.New(log, storage))
		r.Post("/film/add", addFilm.New(log, storage))
		r.Patch("/cassette/available/{cassetteId}", setStatus.New(log, storage))
		r.Patch("/customer/balance/{customerId}", setBalance.New(log, storage))

	})
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", slog.String("error", err.Error()))

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
