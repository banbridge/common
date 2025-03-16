package handler

import (
	"context"

	"log/slog"
)

// MultiHandler 组合多个 slog.Handler 实现多路输出
type MultiHandler []slog.Handler

func (m MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, h := range m {
		if err := h.Handle(ctx, record); err != nil {
			return err
		}
	}
	return nil
}

func (m MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m))
	for i, h := range m {
		handlers[i] = h.WithAttrs(attrs)
	}
	return MultiHandler(handlers)
}

func (m MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m))
	for i, h := range m {
		handlers[i] = h.WithGroup(name)
	}
	return MultiHandler(handlers)
}

var _ slog.Handler = &MultiHandler{}
