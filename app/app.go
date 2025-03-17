package app

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/IceMAN2377/kaspitest/internal/config"
	"github.com/IceMAN2377/kaspitest/internal/repository/postgres"
	"github.com/IceMAN2377/kaspitest/internal/service/egov"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"

	v1Http "github.com/IceMAN2377/kaspitest/internal/transport/http"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

const (
	migrationsPath = "db/migrations"
	dbDriverName   = "postgres"
)

type App struct {
	router *http.ServeMux
	port   int
}

func NewApp(config *config.Config, logger *slog.Logger) *App {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.PostgresHost, config.PostgresPort, config.PostgresUser, config.PostgresPassword, config.PostgresDb, config.PostgresSslMode)

	db, err := sql.Open(dbDriverName, psqlInfo)
	if err != nil {
		panic("failed to connect to db: " + err.Error())
	}

	psql := sqlx.NewDb(db, dbDriverName)

	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.PostgresUser, url.QueryEscape(config.PostgresPassword), config.PostgresHost, config.PostgresPort, config.PostgresDb, config.PostgresSslMode)
	if config.PostgresMigrate {
		if m, err := migrate.New("file://"+migrationsPath, dbURI); err != nil {
			panic("failed to create migrate object: " + err.Error())
		} else if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			panic("failed to apply migrations: " + err.Error())
		}
	}

	repo := postgres.NewRepository(psql)
	service := egov.NewService(repo)
	router := http.NewServeMux()
	v1Http.RegisterEndpoints(logger, router, service)

	return &App{
		router: router,
		port:   config.HttpPort,
	}
}

func (a *App) Run() error {
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", a.port),
		Handler:        a.router,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
