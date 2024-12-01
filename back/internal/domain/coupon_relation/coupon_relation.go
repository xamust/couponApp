package coupon_relation

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

var (
	ErrUserIDCantBeEmpty   = errors.New("user_id cannot be empty")
	ErrCouponIDCantBeEmpty = errors.New("coupon_id cannot be empty")
	ErrCanNotRemove        = errors.New("can't remove applied coupon")
)

type CouponRelation struct {
	id        string
	userID    string
	couponID  string
	metadata  map[string]interface{}
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func (co *CouponRelation) ID() string                       { return co.id }
func (co *CouponRelation) UserID() string                   { return co.userID }
func (co *CouponRelation) CouponID() string                 { return co.couponID }
func (co *CouponRelation) Metadata() map[string]interface{} { return co.metadata }
func (co *CouponRelation) CreatedAt() time.Time             { return co.createdAt }
func (co *CouponRelation) UpdatedAt() time.Time             { return co.updatedAt }
func (co *CouponRelation) DeletedAt() *time.Time            { return co.deletedAt }

type NewCouponRelationInput struct {
	UserID   string
	CouponID string
	Metadata map[string]interface{}
}

func (in *NewCouponRelationInput) Validate() error {
	in.UserID = strings.TrimSpace(in.UserID)
	if in.UserID == "" {
		return ErrUserIDCantBeEmpty
	}
	in.CouponID = strings.TrimSpace(in.CouponID)
	if in.CouponID == "" {
		return ErrCouponIDCantBeEmpty
	}
	if in.Metadata == nil {
		in.Metadata = make(map[string]interface{})
	}
	return nil
}

func New(_ context.Context, in *NewCouponRelationInput) (*CouponRelation, error) {
	if err := in.Validate(); err != nil {
		return nil, fmt.Errorf("validate new coupon relation: %w", err)
	}
	co := &CouponRelation{
		id:       uuid.New().String(),
		userID:   in.UserID,
		couponID: in.CouponID,
		metadata: in.Metadata,
	}
	return co, nil
}

type RestoreSpecification struct {
	ID        string
	UserID    string
	CouponID  string
	Metadata  map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (r *RestoreSpecification) RestoreCouponRelation() *CouponRelation {
	return &CouponRelation{
		id:        r.ID,
		userID:    r.UserID,
		couponID:  r.CouponID,
		metadata:  r.Metadata,
		createdAt: r.CreatedAt,
		updatedAt: r.UpdatedAt,
		deletedAt: r.DeletedAt,
	}
}
