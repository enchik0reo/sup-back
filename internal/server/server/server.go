package server

import (
	"context"
	"net/http"

	"github.com/enchik0reo/sup-back/internal/config"
	"github.com/enchik0reo/sup-back/internal/logs"
)

type Server struct {
	cfg    *config.ApiServer
	log    *logs.CustomLog
	server *http.Server
}

func New(handler http.Handler, c *config.ApiServer, l *logs.CustomLog) *Server {
	srv := setupServer(handler, c)

	return &Server{
		cfg:    c,
		log:    l,
		server: srv,
	}
}

func setupServer(handler http.Handler, cfg *config.ApiServer) *http.Server {
	return &http.Server{
		Addr:           cfg.Address,
		Handler:        handler,
		ReadTimeout:    cfg.Timeout,
		WriteTimeout:   cfg.Timeout,
		IdleTimeout:    cfg.IdleTimeout,
		MaxHeaderBytes: 524288,
	}
}

func (s *Server) Start() error {
	s.log.Info("Web server is running", "address", s.cfg.Address)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
