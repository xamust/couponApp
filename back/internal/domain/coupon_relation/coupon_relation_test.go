package coupon_relation_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xamust/couponApp/internal/domain/coupon_relation"
	"testing"
)

func TestNewAppliedCouponInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   *coupon_relation.NewCouponRelationInput
		wantErr error
	}{
		{
			name: "valid input",
			input: &coupon_relation.NewCouponRelationInput{
				UserID:   "user-123",
				CouponID: "coupon-456",
			},
			wantErr: nil,
		},
		{
			name: "empty userID",
			input: &coupon_relation.NewCouponRelationInput{
				UserID:   "",
				CouponID: "coupon-456",
			},
			wantErr: coupon_relation.ErrUserIDCantBeEmpty,
		},
		{
			name: "empty couponID",
			input: &coupon_relation.NewCouponRelationInput{
				UserID:   "user-123",
				CouponID: "",
			},
			wantErr: coupon_relation.ErrCouponIDCantBeEmpty,
		},
		{
			name: "empty userID and couponID",
			input: &coupon_relation.NewCouponRelationInput{
				UserID:   "",
				CouponID: "",
			},
			wantErr: coupon_relation.ErrUserIDCantBeEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if err != tt.wantErr {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		input      *coupon_relation.NewCouponRelationInput
		wantErr    error
		wantUser   string
		wantCoupon string
	}{
		{
			name: "valid input",
			input: &coupon_relation.NewCouponRelationInput{
				UserID:   "user-123",
				CouponID: "coupon-456",
			},
			wantErr:    nil,
			wantUser:   "user-123",
			wantCoupon: "coupon-456",
		},
		{
			name: "invalid input",
			input: &coupon_relation.NewCouponRelationInput{
				UserID:   "",
				CouponID: "coupon-456",
			},
			wantErr: fmt.Errorf("validate new coupon relation: %w", coupon_relation.ErrUserIDCantBeEmpty),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			couponRelation, err := coupon_relation.New(context.Background(), tt.input)
			if err != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
			if err == nil {
				if couponRelation.UserID() != tt.wantUser {
					t.Errorf("got userID %v, want %v", couponRelation.UserID(), tt.wantUser)
				}
				if couponRelation.CouponID() != tt.wantCoupon {
					t.Errorf("got couponID %v, want %v", couponRelation.CouponID(), tt.wantCoupon)
				}
			}
		})
	}
}
