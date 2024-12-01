package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/xamust/couponApp/internal/domain/user"
)

type UserUsecaseMock struct {
	mock.Mock
}

func (u *UserUsecaseMock) Create(ctx context.Context, in *user.NewUserInput) (*user.User, error) {
	args := u.Called(ctx, in)
	return args.Get(0).(*user.User), args.Error(1)
}

func (u *UserUsecaseMock) GetByID(ctx context.Context, id string) (*user.User, error) {
	args := u.Called(ctx, id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (u *UserUsecaseMock) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	args := u.Called(ctx, limit, offset)
	return args.Get(0).([]*user.User), args.Error(1)
}

func (u *UserUsecaseMock) Delete(ctx context.Context, id string) (*user.User, error) {
	args := u.Called(ctx, id)
	return args.Get(0).(*user.User), args.Error(1)
}
