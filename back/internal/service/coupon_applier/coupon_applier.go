package coupon_applier

import (
	"context"
	"fmt"
	"github.com/xamust/couponApp/internal/domain/coupon"
	"github.com/xamust/couponApp/internal/domain/coupon_relation"
	"github.com/xamust/couponApp/internal/domain/user"
	"golang.org/x/exp/slog"
)

type Repository interface {
	CouponApplier(ctx context.Context, coupon *coupon.Coupon, cr *coupon_relation.CouponRelation) error
}

type CouponApplier interface {
	Applier(ctx context.Context, couponID, userID string) error
}

type couponApplier struct {
	coupRepo    coupon.Repository
	usrRepo     user.Repository
	coupRelRepo coupon_relation.Repository
	applier     Repository
}

func NewCouponApplier(coup coupon.Repository,
	usr user.Repository,
	coupRel coupon_relation.Repository,
	applier Repository) CouponApplier {
	return &couponApplier{
		coupRepo:    coup,
		usrRepo:     usr,
		coupRelRepo: coupRel,
		applier:     applier,
	}
}

func (c *couponApplier) Applier(ctx context.Context, couponID, userID string) error {
	// get user
	one, err := c.usrRepo.FindOne(ctx, userID)
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}

	// get coupon
	oneCoup, err := c.coupRepo.FindOne(ctx, couponID)
	if err != nil {
		return fmt.Errorf("find coupon: %w", err)
	}

	// apply coupon

	// 1. check by already user use coupon
	find, err := c.coupRelRepo.Find(ctx,
		coupon_relation.Cond{CouponIDs: []string{couponID}, UserIDs: []string{userID}},
		coupon_relation.OrderByCreatedDesc,
		0,
		0)
	if err != nil {
		return fmt.Errorf("find coupon relation: %w", err)
	}
	if len(find) > 0 {
		return coupon.ErrAlreadyUsed
	}

	// 2. check by redeemed coupon and max redemptions
	if !oneCoup.IsValid() {
		return coupon.ErrCannotBeApplied
	}

	// 3. create coupon relation
	relation, err := coupon_relation.New(ctx, &coupon_relation.NewCouponRelationInput{
		UserID:   one.ID(),
		CouponID: oneCoup.ID(),
	})
	if err != nil {
		return fmt.Errorf("create coupon relation: %w", err)
	}

	// 4. upd coupon
	if err := oneCoup.Redeem(); err != nil {
		return fmt.Errorf("redeem coupon: %w", err)
	}

	// 5. save coupon relation
	if err := c.applier.CouponApplier(ctx, oneCoup, relation); err != nil {
		return fmt.Errorf("apply coupon: %w", err)
	}

	slog.InfoContext(ctx, "coupon applied for user")
	return nil
}
