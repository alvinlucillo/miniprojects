package models

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

const (
	packageName = "models"
)

func NewPostgreSQLDB(dsn string, logger zerolog.Logger) (*sqlx.DB, error) {
	l := logger.With().Str("package", packageName).Str("function", "NewPostgreSQLDB").Logger()
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		l.Err(err).Msg("failed to open db")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		l.Err(err).Msg("failed to ping db")
		return nil, err
	}
	return db, nil
}
