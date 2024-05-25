package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/enchik0reo/sup-back/internal/config"
	"github.com/enchik0reo/sup-back/internal/logs"
	"github.com/enchik0reo/sup-back/internal/server/handler"
	"github.com/enchik0reo/sup-back/internal/server/server"
	"github.com/enchik0reo/sup-back/internal/storage"
	"github.com/enchik0reo/sup-back/internal/tg"
)

type App struct {
	log   *logs.CustomLog
	cfg   *config.Config
	db    *sql.DB
	tgBot *tg.Bot
	srv   *server.Server
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

	rS := storage.NewReservationStorage(a.db)

	a.tgBot, err = tg.NewBot(rS, a.cfg.TgToken, a.cfg.TgAdmins, a.log)
	if err != nil {
		a.log.Error("Failed to create new tg bot", a.log.Attr("error", err))
		os.Exit(1)
	}

	h := handler.New(a.tgBot, rS, a.cfg.Frontend.Domains, a.cfg.Server.Timeout, a.log)

	a.srv = server.New(h, &a.cfg.Server, a.log)

	return a
}

func (a *App) MustRun() {
	a.log.Info("Starting backend", "env", a.cfg.Env)

	go func() {
		if err := a.srv.Start(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				a.log.Error("Failed over working api sever", a.log.Attr("error", err))
				os.Exit(1)
			}
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		if err := a.tgBot.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				a.log.Error("Failed over working bot", a.log.Attr("error", err))
				return
			}
		}
	}(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	cancel()
	a.mustStop()
}

func (a *App) mustStop() {
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.CtxTimeout)
	defer cancel()

	if err := a.srv.Stop(ctx); err != nil {
		a.log.Error("Closing connection to server", a.log.Attr("error", err))
	}

	if err := a.db.Close(); err != nil {
		a.log.Error("Closing connection to storage", a.log.Attr("error", err))
	}

	a.log.Info("Backend stopped gracefully")
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
