package storage

import (
	"Rest_API/config"
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"log"
	"log/slog"
	"net/url"
	"path/filepath"
)

func GetDBConnect(cfg *config.Config) (*pgx.Conn, error) {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	slog.Info("Connecting to Postgres", "url", dbURL)

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		return nil, errors.Wrap(err, "unable to connect to database")
	}
	log.Println("Running PostgreSQL migrations")
	if err := runPgMigrations(dbURL, "./migrations"); err != nil {
		return nil, errors.Wrap(err, "runPgMigrations failed")
	}
	slog.Info("Connected to PostgreSQL successfully")
	return conn, nil
}
func runPgMigrations(dsn, path string) error {
	if path == "" {
		return errors.New("no migrations path provided")
	}
	if dsn == "" {
		return errors.New("no DSN provided")
	}

	slog.Info("Running migrations...", "dsn", dsn, "path", path)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return errors.Wrap(err, "failed to open DB connection for migrations")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to create migration driver")
	}

	// Convert relative to absolute path and prepend file://
	// ✅ Получаем абсолютный путь с прямыми слэшами

	// 🔥 Преобразуем Windows-путь в корректный file:// URL
	sourceURL, err := getFileMigrationURL("./migrations")
	if err != nil {
		slog.Error("invalid migration path", "error", err)
		return err
	}

	slog.Info("Resolved migration path", "sourceURL", sourceURL)

	m, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		return errors.Wrap(err, "failed to create migrate instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "migration failed")
	}

	slog.Info("Migrations applied successfully")
	return nil
}
func getFileMigrationURL(relPath string) (string, error) {
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return "", err
	}

	// Используем url.PathEscape чтобы защититься от пробелов, кириллицы и т.п.
	u := url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(absPath),
	}
	return u.String(), nil
}
