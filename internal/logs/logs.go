package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	stdLog "log"
	"log/slog"
	"os"
	"runtime"
)

type CustomLog struct {
	*slog.Logger
}

const (
	envLocal     = "local"
	envLocalInfo = "local_info"
	envProd      = "prod"
)

// NewLogger creates a new instance of CustomLog ...
func NewLogger(env string) *CustomLog {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupCustomSlog()
	case envLocalInfo:
		log = setupCustomSlogInfo()
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return &CustomLog{log}
}

// Attr helps to add attributes for logging ...
func (l *CustomLog) Attr(key string, value any) slog.Attr {
	val := slog.AnyValue(value)

	if err, ok := value.(error); ok {
		val = slog.AnyValue(err.Error())
	}

	return slog.Attr{
		Key:   key,
		Value: val,
	}
}

// setupCustomSlog creates a new instance of Logger with logging level debug ...
func setupCustomSlog() *slog.Logger {
	opts := customHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewCustomHandler(os.Stdout)

	return slog.New(handler)
}

// setupCustomSlogInfo creates a new instance of Logger with logging level info ...
func setupCustomSlogInfo() *slog.Logger {
	opts := customHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	}

	handler := opts.NewCustomHandler(os.Stdout)

	return slog.New(handler)
}

type customHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type customHandler struct {
	slog.Handler
	l     *stdLog.Logger
	attrs []slog.Attr
}

// NewCustomHandler is for slog.Handler override ...
func (opts customHandlerOptions) NewCustomHandler(out io.Writer) *customHandler {
	h := &customHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       stdLog.New(out, "", 0),
	}

	return h
}

// NewCustomHandler is for slog.Handler override ...
func (h *customHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[Jan _2 15:04:05]")
	msg := r.Message

	info := runtimeInfo(skipLevel)

	h.l.Println(
		timeStr,
		level,
		msg,
		info,
		string(b),
	)

	return nil
}

// NewCustomHandler is for slog.Handler override ...
func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &customHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

// NewCustomHandler is for slog.Handler override ...
func (h *customHandler) WithGroup(name string) slog.Handler {
	return &customHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}

// Experimental var for runtime info ...
var skipLevel = 4

// runtimeInfo it's xperimental function only for local usage.
// It returns runtime file info ...
func runtimeInfo(skipLevel int) string {
	pc, _, no, ok := runtime.Caller(skipLevel)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return fmt.Sprintf("| %s | line:%d", details.Name(), no)
	}

	return ""
}
