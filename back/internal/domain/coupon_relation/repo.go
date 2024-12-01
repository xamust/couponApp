package coupon_relation

import (
	"context"
)

type OrderBy string

const (
	OrderByCreatedAsc  OrderBy = "created.asc"
	OrderByCreatedDesc OrderBy = "created.desc"
)

type Repository interface {
	Create(ctx context.Context, ca *CouponRelation) error
	FindOne(ctx context.Context, ID string) (*CouponRelation, error)
	Find(ctx context.Context, cond Cond, by OrderBy, limit, offset int) ([]*CouponRelation, error)
	Delete(ctx context.Context, ca *CouponRelation) error
}

type Cond struct {
	Ids       []string
	UserIDs   []string
	CouponIDs []string
	Metadata  map[string]interface{}
}
