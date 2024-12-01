package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/xamust/couponApp/internal/domain/coupon"
)

type CouponUsecaseMock struct {
	mock.Mock
}

func (c *CouponUsecaseMock) Create(ctx context.Context, in *coupon.NewCouponInput) (*coupon.Coupon, error) {
	args := c.Called(ctx, in)
	return args.Get(0).(*coupon.Coupon), args.Error(1)
}

func (c *CouponUsecaseMock) GetByID(ctx context.Context, id string) (*coupon.Coupon, error) {
	args := c.Called(ctx, id)
	return args.Get(0).(*coupon.Coupon), args.Error(1)
}

func (c *CouponUsecaseMock) List(ctx context.Context, limit, offset int) ([]*coupon.Coupon, error) {
	args := c.Called(ctx, limit, offset)
	return args.Get(0).([]*coupon.Coupon), args.Error(1)
}

func (c *CouponUsecaseMock) Delete(ctx context.Context, coup *coupon.Coupon) error {
	args := c.Called(ctx, coup)
	return args.Error(0)
}

func (c *CouponUsecaseMock) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*coupon.Coupon, error) {
	args := c.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*coupon.Coupon), args.Error(1)
}
