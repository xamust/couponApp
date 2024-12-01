package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/xamust/couponApp/internal/domain/coupon"
	"github.com/xamust/couponApp/internal/domain/coupon_relation"
	"github.com/xamust/couponApp/pkg/logger"
	"log/slog"
	"strings"
)

var (
	ErrUserIDEmpty = errors.New("user id is empty")
)

type CouponUsecase interface {
	Create(ctx context.Context, in *coupon.NewCouponInput) (*coupon.Coupon, error)
	GetByID(ctx context.Context, id string) (*coupon.Coupon, error)
	List(ctx context.Context, limit, offset int) ([]*coupon.Coupon, error)
	Delete(ctx context.Context, id string) (*coupon.Coupon, error)
	ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*coupon.Coupon, error)
}

type couponUsecase struct {
	repo         coupon.Repository
	repoRelation coupon_relation.Repository
}

func NewCouponUsecase(
	repo coupon.Repository,
	repoRelation coupon_relation.Repository) CouponUsecase {
	return &couponUsecase{
		repo:         repo,
		repoRelation: repoRelation,
	}
}

func (c *couponUsecase) Create(ctx context.Context, in *coupon.NewCouponInput) (*coupon.Coupon, error) {
	coup, err := coupon.New(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("create coupon: %w", err)
	}
	if err := c.repo.Save(ctx, coup); err != nil {
		return nil, fmt.Errorf("save coupon: %w", err)
	}
	ctx = logger.WithCouponID(ctx, coup.ID())
	slog.DebugContext(ctx, fmt.Sprintf("coupon created: %v\n", coup))
	slog.InfoContext(ctx, "coupon created")
	return coup, nil
}

func (c *couponUsecase) GetByID(ctx context.Context, id string) (*coupon.Coupon, error) {
	one, err := c.repo.FindOne(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find coupon: %w", err)
	}
	slog.InfoContext(ctx, fmt.Sprintf("coupon found: %v\n", one))
	return one, nil
}

func (c *couponUsecase) List(ctx context.Context, limit, offset int) ([]*coupon.Coupon, error) {
	find, err := c.repo.Find(ctx, coupon.Cond{}, coupon.OrderByCreatedDesc, limit, offset)
	if err != nil {
		return nil, err
	}
	slog.DebugContext(ctx, fmt.Sprintf("coupon found: %v\n", find))
	slog.InfoContext(ctx, fmt.Sprintf("%d coupons found\n", len(find)))
	return find, nil
}

func (c *couponUsecase) Delete(ctx context.Context, id string) (*coupon.Coupon, error) {
	coup, err := c.repo.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, coupon.ErrNotFound) {
			slog.InfoContext(ctx, "coupon not found")
			return &coupon.Coupon{}, nil
		}
		return nil, fmt.Errorf("find coupon: %w", err)
	}
	if err := c.repo.Delete(ctx, coup); err != nil {
		return nil, fmt.Errorf("delete coupon: %w", err)
	}
	slog.InfoContext(ctx, fmt.Sprintf("coupon deleted: %v\n", coup))
	return coup, nil
}

func (c *couponUsecase) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*coupon.Coupon, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrUserIDEmpty
	}
	find, err := c.repoRelation.Find(ctx, coupon_relation.Cond{
		UserIDs: []string{userID},
	}, coupon_relation.OrderByCreatedDesc, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(find) == 0 {
		slog.InfoContext(ctx, "coupon not found")
		return nil, nil
	}
	couponIDs := make([]string, len(find))
	for i, fr := range find {
		couponIDs[i] = fr.CouponID()
	}
	coupons, err := c.repo.Find(ctx, coupon.Cond{
		Ids: couponIDs,
	}, coupon.OrderByCreatedDesc, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("find coupon: %w", err)
	}
	slog.DebugContext(ctx, fmt.Sprintf("coupons found: %v\n", find))
	slog.InfoContext(ctx, fmt.Sprintf("%d coupons found\n", len(coupons)))
	return coupons, nil
}
