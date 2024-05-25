package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/enchik0reo/sup-back/internal/logs"
	"github.com/enchik0reo/sup-back/internal/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Notifier interface {
	PushNotice() error
}

type Storage interface {
	GetReserved(ctx context.Context, from, to string) ([]models.Sup, error)
	CreateApprove(ctx context.Context, approve models.Approve) (int64, error)
}

type CustomRouter struct {
	*chi.Mux

	notifier Notifier
	storage  Storage

	timeout time.Duration
	log     *logs.CustomLog
}

func New(n Notifier, s Storage, domains []string, timeout time.Duration, log *logs.CustomLog) http.Handler {
	r := CustomRouter{chi.NewRouter(), n, s, timeout, log}

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(loggerMw(log))
	r.Use(corsSettings(domains))

	r.Post("/api/v1/getItems", r.getItems())
	r.Post("/api/v1/makeReservation", r.makeReservation())

	return r
}
