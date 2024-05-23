package app

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/enchik0reo/sup-back/internal/config"
	"github.com/enchik0reo/sup-back/internal/logs"
	"github.com/enchik0reo/sup-back/internal/storage"
)

type App struct {
	log *logs.CustomLog
	cfg *config.Config
	db  *sql.DB
}

func New() *App {
	var err error
	a := &App{}

	a.cfg = config.MustLoad()

	a.log = logs.NewLogger(a.cfg.Env)

	a.db, err = connectionAttemptToDB(a.cfg.Storage)
	if err != nil {
		a.log.Error("Failed to connect to db", a.log.Attr("error", err))
		os.Exit(1)
	}

	return a
}

func (a *App) MustRun() {

}

func (a *App) mustStop() {

}

func connectionAttemptToDB(psql config.Postgres) (*sql.DB, error) {
	var err error
	var db *sql.DB

	for i := 1; i <= 5; i++ {
		db, err = storage.Connect(psql)
		if err != nil {
			time.Sleep(time.Duration(i) * time.Second)
		} else {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("after retries: %w", err)
	}

	return db, nil
}
