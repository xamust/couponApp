package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/xamust/couponApp/internal/domain/coupon"
)

type RepoMock struct {
	mock.Mock
}

func (m *RepoMock) FindOne(ctx context.Context, id string) (*coupon.Coupon, error) {
	args := m.Called(id)
	return args.Get(0).(*coupon.Coupon), args.Error(1)
}

func (m *RepoMock) Find(ctx context.Context, cond coupon.Cond, by coupon.OrderBy, limit, offset int) ([]*coupon.Coupon, error) {
	args := m.Called(cond, by, limit, offset)
	return args.Get(0).([]*coupon.Coupon), args.Error(1)
}

func (m *RepoMock) Save(ctx context.Context, co *coupon.Coupon) error {
	args := m.Called(co)
	return args.Error(0)
}

func (m *RepoMock) Delete(ctx context.Context, id string) (*coupon.Coupon, error) {
	args := m.Called(id)
	return args.Get(0).(*coupon.Coupon), args.Error(1)
}
