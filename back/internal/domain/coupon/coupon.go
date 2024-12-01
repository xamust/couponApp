package coupon

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"strings"
	"time"
)

var (
	ErrIncorrectRedemBy  = errors.New("incorrect redeem by")
	ErrNameCantBeEmpty   = errors.New("name cannot be empty")
	ErrRewardCantBeEmpty = errors.New("reward cannot be empty")
	ErrNotValid          = errors.New("coupon not valid")
	ErrCanNotRemove      = errors.New("can't remove coupon")
	ErrAlreadyUsed       = errors.New("coupon already used")
	ErrCannotBeApplied   = errors.New("coupon cannot be applied")
)

type Coupon struct {
	id     string
	name   string
	reward string
	// maxRedemptions лимиты на использование купона (если "0" - неограниченное использование)
	maxRedemptions int

	// timesRedeemed количество использований купона
	timesRedeemed int

	// redeemBy купон нельзя будет применить после этого времени (если nil - неограниченное время)
	redeemBy *time.Time

	metadata  map[string]string
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func (co *Coupon) ID() string                  { return co.id }
func (co *Coupon) Name() string                { return co.name }
func (co *Coupon) Reward() string              { return co.reward }
func (co *Coupon) MaxRedemptions() int         { return co.maxRedemptions }
func (co *Coupon) TimesRedeemed() int          { return co.timesRedeemed }
func (co *Coupon) RedeemBy() *time.Time        { return co.redeemBy }
func (co *Coupon) Metadata() map[string]string { return co.metadata }
func (co *Coupon) CreatedAt() time.Time        { return co.createdAt }
func (co *Coupon) UpdatedAt() time.Time        { return co.updatedAt }
func (co *Coupon) DeletedAt() *time.Time       { return co.deletedAt }

// IsValid разрешен ли купон к использованию
func (co *Coupon) IsValid() bool {
	if co.maxRedemptions > 0 && co.maxRedemptions <= co.timesRedeemed {
		return false
	}
	if co.redeemBy != nil {
		return time.Now().Before(*co.redeemBy)
	}
	return true
}

type NewCouponInput struct {
	Name           string `validate:"gte=3,lte=30"`
	Reward         string `validate:"gte=3,lte=255"`
	MaxRedemptions int
	RedeemBy       *time.Time
	Metadata       map[string]string
}

func (in *NewCouponInput) Validate() error {
	validate := validator.New()
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return ErrNameCantBeEmpty
	}

	in.Reward = strings.TrimSpace(in.Reward)
	if in.Reward == "" {
		return ErrRewardCantBeEmpty
	}

	if in.RedeemBy != nil {
		if time.Now().After(*in.RedeemBy) {
			return ErrIncorrectRedemBy
		}
	}
	if in.Metadata == nil {
		in.Metadata = make(map[string]string)
	}
	return validate.Struct(in)
}

func New(_ context.Context, in *NewCouponInput) (*Coupon, error) {
	if err := in.Validate(); err != nil {
		return nil, fmt.Errorf("validate new coupon: %w", err)
	}

	co := &Coupon{
		id:             uuid.New().String(),
		name:           in.Name,
		reward:         in.Reward,
		maxRedemptions: in.MaxRedemptions,
		timesRedeemed:  0,
		redeemBy:       in.RedeemBy,
		metadata:       in.Metadata,
	}
	return co, nil
}

func (co *Coupon) Redeem() error {
	if !co.IsValid() {
		return ErrNotValid
	}
	co.timesRedeemed++
	return nil
}

type RestoreSpecification struct {
	ID             string
	Name           string
	Reward         string
	MaxRedemptions int
	TimesRedeemed  int
	RedeemBy       *time.Time
	Metadata       map[string]string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

func (r *RestoreSpecification) RestoreCoupon() *Coupon {
	return &Coupon{
		id:             r.ID,
		name:           r.Name,
		reward:         r.Reward,
		maxRedemptions: r.MaxRedemptions,
		timesRedeemed:  r.TimesRedeemed,
		redeemBy:       r.RedeemBy,
		metadata:       r.Metadata,
		createdAt:      r.CreatedAt,
		updatedAt:      r.UpdatedAt,
		deletedAt:      r.DeletedAt,
	}
}
