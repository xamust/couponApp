package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/xamust/couponApp/internal/domain/coupon_relation"
)

type RepoMock struct {
	mock.Mock
}

func (r *RepoMock) Create(ctx context.Context, couponID string, userID string) error {
	args := r.Called(couponID, userID)
	return args.Error(0)
}

func (r *RepoMock) FindByID(ctx context.Context, ID string) (coupon_relation.CouponRelation, error) {
	args := r.Called(ID)
	return args.Get(0).(coupon_relation.CouponRelation), args.Error(1)
}

func (r *RepoMock) FindByUserID(ctx context.Context, userID string) ([]coupon_relation.CouponRelation, error) {
	args := r.Called(userID)
	return args.Get(0).([]coupon_relation.CouponRelation), args.Error(1)
}

func (r *RepoMock) FindByCouponID(ctx context.Context, couponID string) ([]coupon_relation.CouponRelation, error) {
	args := r.Called(couponID)
	return args.Get(0).([]coupon_relation.CouponRelation), args.Error(1)
}

func (r *RepoMock) DeleteByID(ctx context.Context, ID string) error {
	args := r.Called(ID)
	return args.Error(0)
}
