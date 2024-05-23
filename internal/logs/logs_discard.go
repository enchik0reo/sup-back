package logs

import (
	"context"
	"io"
	"log/slog"
)

type discardHandler struct {
	out io.Writer
}

// Discard logger for tests
func NewDiscardLogger() *CustomLog {
	log := slog.New(newDiscardHandler(io.Discard))

	return &CustomLog{log}
}

func newDiscardHandler(out io.Writer) *discardHandler {
	return &discardHandler{out: out}
}

func (h *discardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *discardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *discardHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *discardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
