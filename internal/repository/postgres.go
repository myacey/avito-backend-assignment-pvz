package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/myacey/avito-backend-assignment-pvz/internal/config"
	db "github.com/myacey/avito-backend-assignment-pvz/internal/repository/sqlc"
)

func ConfigurePostgres(cfg config.AppConfig) (*db.Queries, *sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB,
	)

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open postgres conn: %w", err)
	}
	err = conn.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot ping postgres: %w", err)
	}

	conn.SetMaxOpenConns(500)
	conn.SetMaxIdleConns(400)
	conn.SetConnMaxIdleTime(5 * time.Minute)

	queries := db.New(conn)

	return queries, conn, nil
}

// isUniqueViolation checks if err is about
// not unique val.
func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
