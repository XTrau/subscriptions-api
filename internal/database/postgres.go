package database

import (
	"database/sql"
	"fmt"
	"subscriptions-api/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func RunMigrations(cfg config.Config) error {
	m, err := migrate.New(
		"file://migrations",
		GetPostgresDsn(cfg),
	)

	if err != nil {
		return err
	}

	defer m.Close()

	err = m.Up()
	return err
}

func GetPostgresDsn(cfg config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
}

func NewPostresDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", GetPostgresDsn(cfg))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
