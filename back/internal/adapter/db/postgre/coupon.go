package postgre

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/xamust/couponApp/internal/domain/coupon"
	"github.com/xamust/couponApp/utils/wrappers"
	"gorm.io/gorm"
	"strings"
	"time"
)

type CouponRepository struct {
	db_ *gorm.DB
}

func NewCouponRepository(db *gorm.DB) coupon.Repository {
	return &CouponRepository{db_: db}
}

func (rep *CouponRepository) db(ctx context.Context) *gorm.DB {
	return rep.db_.WithContext(ctx)
}

func (rep *CouponRepository) FindOne(ctx context.Context, id string) (*coupon.Coupon, error) {
	co := &Coupon{}
	q := rep.applyCond(rep.db(ctx), coupon.Cond{Ids: []string{strings.ToLower(id)}})

	if err := q.First(co).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, coupon.ErrNotFound
		}
		return nil, err
	}
	return UnmarshalCoupon(co)
}

func (rep *CouponRepository) Find(ctx context.Context, cond coupon.Cond, by coupon.OrderBy, limit, offset int) ([]*coupon.Coupon, error) {
	models := make([]*Coupon, 0)

	q := rep.applyCond(rep.db(ctx), cond)
	q = rep.applyOrder(q, by)
	q = rep.applyLimit(q, limit, offset)

	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}

	coupons := make([]*coupon.Coupon, len(models))
	for k, v := range models {
		coup, err := UnmarshalCoupon(v)
		if err != nil {
			return nil, err
		}
		coupons[k] = coup
	}
	return coupons, nil
}

func (rep *CouponRepository) applyCond(q *gorm.DB, cond coupon.Cond) *gorm.DB {
	q = q.Table("coupons")

	if len(cond.Ids) > 0 {
		q = q.Where("coupons.id IN (?)", cond.Ids)
	}
	if cond.RedeemBy != nil {
		q = q.Where("coupons.redeem_by <= ?)", cond.RedeemBy)
	}
	if len(cond.Metadata) > 0 {
		for k, v := range cond.Metadata {
			q = q.Where("coupons.metadata->>? IN (?)", k, v)
		}
	}
	return q
}

func (rep *CouponRepository) applyOrder(q *gorm.DB, by coupon.OrderBy) *gorm.DB {
	switch by {
	case coupon.OrderByCreatedAsc:
		q = q.Order("coupons.created_at asc")
	default:
		q = q.Order("coupons.created_at desc")
	}
	return q
}

func (rep *CouponRepository) applyLimit(q *gorm.DB, limit, offset int) *gorm.DB {
	if limit > 0 {
		q = q.Limit(limit)
	}
	if limit > 100 {
		q = q.Limit(100)
	}
	return q.Offset(offset)
}

func (rep *CouponRepository) Save(ctx context.Context, co *coupon.Coupon) error {
	marshalCoupon, err := MarshalCoupon(co)
	if err != nil {
		return err
	}
	if err := rep.db(ctx).Save(marshalCoupon).Error; err != nil {
		return err
	}
	return nil
}

func (rep *CouponRepository) Delete(ctx context.Context, coup *coupon.Coupon) error {
	if err := rep.db_.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		marshalCoupon, err := MarshalCoupon(coup)
		if err != nil {
			tx.Rollback()
			return err
		}
		callback := tx.Model(marshalCoupon).Where("times_redeemed = 0").Update("deleted_at", time.Now())
		if callback.Error != nil {
			return callback.Error
		}
		if callback.RowsAffected == 0 {
			tx.Rollback()
			return coupon.ErrCanNotRemove
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

type Coupon struct {
	ID             string `gorm:"PrimaryKey"`
	Name           string
	Reward         string
	MaxRedemptions int
	TimesRedeemed  int
	RedeemBy       *time.Time
	Metadata       []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
}

func (Coupon) TableName() string { return "coupons" }

func MarshalCoupon(co *coupon.Coupon) (*Coupon, error) {
	md, err := json.Marshal(co.Metadata())
	if err != nil {
		return nil, err
	}
	return &Coupon{
		ID:             co.ID(),
		Name:           co.Name(),
		Reward:         co.Reward(),
		MaxRedemptions: co.MaxRedemptions(),
		TimesRedeemed:  co.TimesRedeemed(),
		RedeemBy:       co.RedeemBy(),
		Metadata:       md,
		CreatedAt:      co.CreatedAt(),
		UpdatedAt:      co.UpdatedAt(),
		DeletedAt:      wrappers.WrapGormDeletedAt(co.DeletedAt()),
	}, nil
}

func UnmarshalCoupon(co *Coupon) (*coupon.Coupon, error) {
	md := map[string]string{}
	if co.Metadata != nil {
		if err := json.Unmarshal(co.Metadata, &md); err != nil {
			return nil, err
		}
	}
	restoreCoupon := coupon.RestoreSpecification{
		ID:             co.ID,
		Name:           co.Name,
		Reward:         co.Reward,
		MaxRedemptions: co.MaxRedemptions,
		TimesRedeemed:  co.TimesRedeemed,
		RedeemBy:       co.RedeemBy,
		Metadata:       md,
		CreatedAt:      co.CreatedAt,
		UpdatedAt:      co.UpdatedAt,
		DeletedAt:      wrappers.UnwrapGormDeletedAt(co.DeletedAt),
	}
	return restoreCoupon.RestoreCoupon(), nil
}
