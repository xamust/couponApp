package logger

import (
	"context"
)

type logCtx struct {
	CouponID   string
	CouponName string

	UserID   string
	UserName string
}
type keyType int

const key = keyType(0)

func WithCouponID(ctx context.Context, couponID string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.CouponID = couponID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{CouponID: couponID})
}

func WithCouponName(ctx context.Context, couponName string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.CouponName = couponName
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{CouponName: couponName})
}

func WithUserID(ctx context.Context, userID string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.UserID = userID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{UserID: userID})
}

func WithUserName(ctx context.Context, userName string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.UserName = userName
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{UserName: userName})
}
