package logger

import (
	"log/slog"
)

type logger struct {
	slog *slog.Logger
	opts Options
}

func NewLogger(o ...Option) {
	l := &logger{
		opts: NewOptions(o...),
	}
	handler := slog.Handler(slog.NewJSONHandler(l.opts.output, &slog.HandlerOptions{Level: l.opts.level}))
	handler = NewHandlerMiddleware(handler)

	slog.SetDefault(slog.New(handler))
}
