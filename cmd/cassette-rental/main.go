package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	orderscassettesrepo "github.com/pancher1990/cassette-rental/internal/repositories/orders-cassettes"
	"github.com/pancher1990/cassette-rental/internal/usecases/cassettes"
	"github.com/pancher1990/cassette-rental/internal/usecases/orders"
	"github.com/pancher1990/cassette-rental/internal/usecases/sessions"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/pancher1990/cassette-rental/internal/controllers/api"
	cassettesrepo "github.com/pancher1990/cassette-rental/internal/repositories/cassettes"
	customersrepo "github.com/pancher1990/cassette-rental/internal/repositories/customers"
	filmsrepo "github.com/pancher1990/cassette-rental/internal/repositories/films"
	orsersrepo "github.com/pancher1990/cassette-rental/internal/repositories/orders"
	sessionsrepo "github.com/pancher1990/cassette-rental/internal/repositories/sessions"

	"github.com/pancher1990/cassette-rental/internal/transaction"
	"github.com/pancher1990/cassette-rental/internal/usecases/customers"
	"github.com/pancher1990/cassette-rental/internal/usecases/films"
)

type PostgresConfig struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     uint16 `envconfig:"PORT" default:"5432"`
	Name     string `envconfig:"NAME" required:"true"`
	User     string `envconfig:"USER" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	SSLMode  string `envconfig:"SSL_MODE" default:"disable"`
}

func (p PostgresConfig) dsn() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.Name,
		p.SSLMode,
	)
}

type config struct {
	Database   PostgresConfig `envconfig:"DATABASE"`
	HTTPServer struct {
		Host        string        `envconfig:"HOST"`
		Port        uint16        `envconfig:"PORT" default:"8080"`
		Timeout     time.Duration `envconfig:"TIMEOUT" default:"4s"`
		IdleTimeout time.Duration `envconfig:"IDLE_TIMEOUT" default:"60s"`
	} `envconfig:"HTTP_SERVER"`
	Admin struct {
		Login    string `envconfig:"LOGIN" default:"admin"`
		Password string `envconfig:"PASSWORD" required:"true"`
	} `envconfig:"ADMIN"`
}

func main() {
	logger := newLogger()

	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		logger.Error("failed to get config", slog.String("err", err.Error()))

		return
	}

	customerRepo := customersrepo.New()
	filmRepo := filmsrepo.New()
	orderRepo := orsersrepo.New()
	cassettesRepo := cassettesrepo.New()
	ordersCassettesRepo := orderscassettesrepo.New()
	sessionsRepo := sessionsrepo.New()

	pool, err := newPgxPool(cfg.Database.dsn())
	if err != nil {
		logger.Error("failed to create database pool", slog.String("err", err.Error()))

		return
	}

	controller, err := api.New(
		api.WithCustomerCreater(customers.Create(customerRepo, transaction.Tx(pool, logger))),
		api.WithCustomerBalanceUpdater(customers.UpdateBalance(customerRepo, transaction.Tx(pool, logger))),
		api.WithFilmCreater(films.Create(filmRepo, transaction.Tx(pool, logger))),
		api.WithFilmFinder(films.Find(filmRepo, transaction.Tx(pool, logger))),
		api.WithOrderCreater(orders.Create(orders.Repositories{
			Film:          filmRepo,
			Customer:      customerRepo,
			Cassette:      cassettesRepo,
			Order:         orderRepo,
			OrderCassette: ordersCassettesRepo,
		}, transaction.Tx(pool, logger))),
		api.WithCassetteCreater(cassettes.Create(cassettes.Repositories{
			Film:     filmRepo,
			Cassette: cassettesRepo,
		}, transaction.Tx(pool, logger))),
		api.WithLogin(sessions.Login(sessions.Repositories{
			SessionRepository:  sessionsRepo,
			CustomerRepository: customerRepo,
		}, transaction.Tx(pool, logger))),
		api.WithLogout(sessions.Logout(sessionsRepo, transaction.Tx(pool, logger))),
	)
	if err != nil {
		logger.Error("failed to create controller", slog.String("err", err.Error()))

		return
	}

	server := http.Server{
		Addr:        fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port),
		ReadTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
		Handler:     controller,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("failed to start server", slog.String("err", err.Error()))

		return
	}
}

func newLogger() *slog.Logger {
	level := slog.LevelInfo
	if err := level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))); err != nil {
		level = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}

func newPgxPool(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	return pool, nil
}
