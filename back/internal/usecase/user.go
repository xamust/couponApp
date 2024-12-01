package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/xamust/couponApp/internal/domain/user"
	"github.com/xamust/couponApp/pkg/logger"
	"log/slog"
)

type UserUsecase interface {
	Create(ctx context.Context, in *user.NewUserInput) (*user.User, error)
	GetByID(ctx context.Context, id string) (*user.User, error)
	List(ctx context.Context, limit, offset int) ([]*user.User, error)
	Delete(ctx context.Context, id string) (*user.User, error)
}

type userUsecase struct {
	repo user.Repository
}

func NewUserUsecase(repo user.Repository) UserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

func (u *userUsecase) Create(ctx context.Context, in *user.NewUserInput) (*user.User, error) {
	usr, err := user.New(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	if err := u.repo.Save(ctx, usr); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}
	ctx = logger.WithUserID(ctx, usr.ID())
	slog.InfoContext(ctx, fmt.Sprintf("user created: %v\n", usr))
	return usr, nil
}

func (u *userUsecase) GetByID(ctx context.Context, id string) (*user.User, error) {
	one, err := u.repo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	slog.InfoContext(ctx, fmt.Sprintf("user found: %v\n", one))
	return one, nil
}

func (u *userUsecase) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	find, err := u.repo.Find(ctx, user.Cond{}, user.OrderByCreatedDesc, limit, offset)
	if err != nil {
		return nil, err
	}
	slog.DebugContext(ctx, fmt.Sprintf("users found: %v", find))
	slog.InfoContext(ctx, fmt.Sprintf("%d users found", len(find)))
	return find, nil
}

func (u *userUsecase) Delete(ctx context.Context, id string) (*user.User, error) {
	usr, err := u.repo.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			slog.InfoContext(ctx, "user not found")
			return &user.User{}, nil
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
	if err := u.repo.Delete(ctx, usr); err != nil {
		return nil, fmt.Errorf("delete user: %w", err)
	}
	slog.InfoContext(ctx, fmt.Sprintf("user deleted: %v", usr))
	return usr, nil
}
