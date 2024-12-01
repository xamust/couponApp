package coupon

import (
	"context"
	"errors"
	"time"
)

type OrderBy string

const (
	OrderByCreatedAsc  OrderBy = "created.asc"
	OrderByCreatedDesc OrderBy = "created.desc"
)

var ErrNotFound = errors.New("coupon not found")

type Repository interface {
	FindOne(ctx context.Context, id string) (*Coupon, error)
	Find(ctx context.Context, cond Cond, by OrderBy, limit, offset int) ([]*Coupon, error)
	Save(ctx context.Context, co *Coupon) error
	Delete(ctx context.Context, coup *Coupon) error
}

type Cond struct {
	Ids      []string
	RedeemBy *time.Time
	Metadata map[string]interface{}
}
