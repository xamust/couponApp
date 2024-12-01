package logger

import (
	"context"
	"log/slog"
)

type HandlerMiddlware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddlware {
	return &HandlerMiddlware{next: next}
}
func (h *HandlerMiddlware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}
func (h *HandlerMiddlware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(key).(logCtx); ok {
		if c.CouponID != "" {
			rec.Add("coupon_id", c.CouponID)
		}
		if c.CouponName != "" {
			rec.Add("coupon_name", c.CouponName)
		}
		if c.UserID != "" {
			rec.Add("user_id", c.UserID)
		}
		if c.UserName != "" {
			rec.Add("user_name", c.UserName)
		}
	}
	return h.next.Handle(ctx, rec)
}
func (h *HandlerMiddlware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithAttrs(attrs)} // todo @ не забыть обернуть!!!
}
func (h *HandlerMiddlware) WithGroup(name string) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithGroup(name)} // todo @ не забыть обернуть!!!
}
