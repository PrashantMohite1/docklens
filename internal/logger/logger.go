package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type Handler struct {
	w io.Writer
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	fmt.Fprintln(h.w, r.Message)
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return h
}

func Init() {
	handler := &Handler{
		w: os.Stdout,
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
