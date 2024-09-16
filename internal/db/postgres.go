package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"zadanie-6105/internal/config"

	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	var connStr string

	if cfg.PostgresConn != "" {
		connStr = cfg.PostgresConn
	} else {
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.PostgresHost,
			strconv.Itoa(cfg.PostgresPort),
			cfg.PostgresUsername,
			cfg.PostgresPassword,
			cfg.PostgresDatabase,
		)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL database!")
	return db, nil
}
