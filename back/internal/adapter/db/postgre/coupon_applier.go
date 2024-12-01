package postgre

import (
	"context"
	"fmt"
	"github.com/xamust/couponApp/internal/domain/coupon"
	"github.com/xamust/couponApp/internal/domain/coupon_relation"
	"github.com/xamust/couponApp/internal/service/coupon_applier"
	"gorm.io/gorm"
)

type CouponApplierRepository struct {
	db_ *gorm.DB
}

func NewCouponApplierRepository(db *gorm.DB) coupon_applier.Repository {
	return &CouponApplierRepository{db_: db}
}

func (rep *CouponApplierRepository) db(ctx context.Context) *gorm.DB {
	return rep.db_.WithContext(ctx)
}

func (rep *CouponApplierRepository) CouponApplier(ctx context.Context, coupon *coupon.Coupon, cr *coupon_relation.CouponRelation) error {
	if coupon == nil {
		return fmt.Errorf("coupon is nil")
	}
	if cr == nil {
		return fmt.Errorf("coupon relation is nil")
	}
	return rep.db_.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		coupRepo := NewCouponRepository(tx)
		couRelationRepo := NewCouponRelationRepository(tx)

		if err := coupRepo.Save(ctx, coupon); err != nil {
			return fmt.Errorf("save coupon: %w", err)
		}
		if err := couRelationRepo.Create(ctx, cr); err != nil {
			return fmt.Errorf("save coupon relation: %w", err)
		}
		return nil
	})
}
