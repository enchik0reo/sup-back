package storage

import (
	"database/sql"
	"fmt"

	"github.com/enchik0reo/sup-back/internal/config"

	_ "github.com/lib/pq"
)

func Connect(cfg config.Postgres) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Driver, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
