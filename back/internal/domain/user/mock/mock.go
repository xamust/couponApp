package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/xamust/couponApp/internal/domain/user"
)

type RepoMock struct {
	mock.Mock
}

func (m *RepoMock) FindOne(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *RepoMock) Find(ctx context.Context, cond user.Cond, by user.OrderBy, limit, offset int) ([]*user.User, error) {
	args := m.Called(cond, by, limit, offset)
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *RepoMock) Save(ctx context.Context, u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *RepoMock) Delete(ctx context.Context, u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}
