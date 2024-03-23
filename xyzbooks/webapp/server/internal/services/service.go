package service

import (
	"alvinlucillo/xyzbooks_webapp/internal/models"
	"alvinlucillo/xyzbooks_webapp/internal/repos"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

const (
	DBTypePostgres = "postgres"
	packageName    = "service"
)

type Service struct {
	Repository repos.Repository
	tx         *sqlx.Tx
	db         *sqlx.DB
	logger     zerolog.Logger
}

type ServiceConfig struct {
	DBHost     string
	DBPort     string
	DBPassword string
	DBUser     string
	DBName     string
	DBType     string
	Logger     zerolog.Logger
}

func NewService(cfg ServiceConfig) (*Service, error) {
	l := cfg.Logger.With().Str("package", packageName).Str("function", "NewService").Logger()

	// Create dsn
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Create postgres db (since we only have postgres for now)
	db, err := models.NewPostgreSQLDB(dsn, cfg.Logger)
	if err != nil {
		l.Err(err).Msg("failed to create postgres db")
		return nil, err
	}

	return &Service{
		db:     db,
		logger: cfg.Logger,
	}, nil
}

func (svc *Service) InitRepo() error {
	l := svc.logger.With().Str("package", packageName).Str("function", "InitRepo").Logger()

	tx, err := svc.db.Beginx()
	if err != nil {
		l.Err(err).Msg("failed to begin transaction")
		return err
	}

	repo := repos.NewRepository(tx, svc.logger)

	svc.tx = tx
	svc.Repository = repo

	return nil
}

func (svc *Service) SendError(w http.ResponseWriter, httpCode int, errVal interface{}) {
	l := svc.logger.With().Str("package", packageName).Str("function", "SendError").Logger()

	err := svc.tx.Rollback()
	if err != nil {
		l.Err(err).Msg("failed to rollback transaction")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)

	switch httpCode {
	case http.StatusInternalServerError:
		json.NewEncoder(w).Encode("Internal server error")
	case http.StatusNotFound:
		json.NewEncoder(w).Encode("Record not found")
	default:
		json.NewEncoder(w).Encode(errVal)
	}
}

func (svc *Service) SendResponse(w http.ResponseWriter, request *http.Request, data interface{}) {
	l := svc.logger.With().Str("package", packageName).Str("function", "SendResponse").Logger()
	err := svc.tx.Commit()
	if err != nil {
		l.Err(err).Msg("failed to commit transaction")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
