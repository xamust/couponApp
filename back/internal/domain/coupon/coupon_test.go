package coupon_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/xamust/couponApp/internal/domain/coupon"
	"testing"
	"time"
)

var (
	now    = time.Now()
	past   = now.Add(-24 * time.Hour)
	future = now.Add(24 * time.Hour)
)

func TestNewCouponInput_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     *coupon.NewCouponInput
		expectErr error
	}{
		{
			name: "valid input",
			input: &coupon.NewCouponInput{
				Name:     "Test Coupon",
				Reward:   "Free Coffee",
				RedeemBy: &future,
			},
			expectErr: nil,
		},
		{
			name: "empty name",
			input: &coupon.NewCouponInput{
				Name:     "",
				Reward:   "Free Coffee",
				RedeemBy: &future,
			},
			expectErr: coupon.ErrNameCantBeEmpty,
		},
		{
			name: "empty redemBy",
			input: &coupon.NewCouponInput{
				Name:   "Coffee Coupon",
				Reward: "Free Coffee",
			},
			expectErr: nil,
		},
		{
			name: "empty reward",
			input: &coupon.NewCouponInput{
				Name:     "Test Coupon",
				Reward:   "",
				RedeemBy: &future,
			},
			expectErr: coupon.ErrRewardCantBeEmpty,
		},
		{
			name: "redeemBy in the past",
			input: &coupon.NewCouponInput{
				Name:     "Test Coupon",
				Reward:   "Free Coffee",
				RedeemBy: &past,
			},
			expectErr: coupon.ErrIncorrectRedemBy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			assert.Equal(t, tt.expectErr, err)
		})
	}
}

func TestCoupon_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		coupon *coupon.RestoreSpecification
		valid  bool
	}{
		{
			name: "valid coupon with no limits",
			coupon: &coupon.RestoreSpecification{
				MaxRedemptions: 0,
				TimesRedeemed:  0,
				RedeemBy:       nil,
			},
			valid: true,
		},
		{
			name: "coupon expired by redeemBy",
			coupon: &coupon.RestoreSpecification{
				MaxRedemptions: 0,
				TimesRedeemed:  0,
				RedeemBy:       &past,
			},
			valid: false,
		},
		{
			name: "coupon with max redemptions reached",
			coupon: &coupon.RestoreSpecification{
				MaxRedemptions: 1,
				TimesRedeemed:  1,
				RedeemBy:       &future,
			},
			valid: false,
		},
		{
			name: "valid coupon with limits",
			coupon: &coupon.RestoreSpecification{
				MaxRedemptions: 2,
				TimesRedeemed:  1,
				RedeemBy:       &future,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.coupon.RestoreCoupon().IsValid())
		})
	}
}

func TestCoupon_Redeem(t *testing.T) {
	tests := []struct {
		name      string
		coupon    *coupon.RestoreSpecification
		expectErr error
	}{
		{
			name: "successful redeem",
			coupon: &coupon.RestoreSpecification{
				MaxRedemptions: 2,
				TimesRedeemed:  1,
				RedeemBy:       &future,
			},
			expectErr: nil,
		},
		{
			name: "redeem expired coupon",
			coupon: &coupon.RestoreSpecification{
				MaxRedemptions: 0,
				TimesRedeemed:  0,
				RedeemBy:       &past,
			},
			expectErr: coupon.ErrNotValid,
		},
		{
			name: "redeem max redemptions reached",
			coupon: &coupon.RestoreSpecification{
				MaxRedemptions: 1,
				TimesRedeemed:  1,
				RedeemBy:       &future,
			},
			expectErr: coupon.ErrNotValid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.coupon.RestoreCoupon().Redeem()
			assert.Equal(t, tt.expectErr, err)
		})
	}
}

func TestNewCoupon(t *testing.T) {
	tests := []struct {
		name      string
		input     *coupon.NewCouponInput
		expectErr error
	}{
		{
			name: "create valid coupon",
			input: &coupon.NewCouponInput{
				Name:           "Test Coupon",
				Reward:         "Discount",
				MaxRedemptions: 10,
				RedeemBy:       &future,
				Metadata:       map[string]string{"category": "promo"},
			},
			expectErr: nil,
		},
		{
			name: "invalid input - empty name",
			input: &coupon.NewCouponInput{
				Name:           "",
				Reward:         "Discount",
				MaxRedemptions: 10,
				RedeemBy:       &future,
				Metadata:       map[string]string{"category": "promo"},
			},
			expectErr: coupon.ErrNameCantBeEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coupon, err := coupon.New(context.Background(), tt.input)
			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Name, coupon.Name())
				assert.Equal(t, tt.input.Reward, coupon.Reward())
			}
		})
	}
}

func TestUser_Methods(t *testing.T) {
	redemBy := time.Now().Add(1 * time.Hour)
	coup := &coupon.RestoreSpecification{
		ID:             "123",
		Name:           "Test Coupon",
		Reward:         "Free Coffee",
		MaxRedemptions: 0,
		TimesRedeemed:  10,
		RedeemBy:       &redemBy,
		Metadata:       map[string]string{"test": "test"},
		CreatedAt:      now,
		UpdatedAt:      now,
		DeletedAt:      nil,
	}
	coupon := coup.RestoreCoupon()
	assert.Equal(t, "123", coupon.ID(), "ID() should return correct ID")
	assert.Equal(t, "Test Coupon", coupon.Name(), "Name() should return correct name")
	assert.Equal(t, "Free Coffee", coupon.Reward(), "Reward() should return correct reward")
	assert.Equal(t, 0, coupon.MaxRedemptions(), "MaxRedemptions() should return correct max redemptions")
	assert.Equal(t, 10, coupon.TimesRedeemed(), "TimesRedeemed() should return correct times redeemed")
	assert.Equal(t, &redemBy, coupon.RedeemBy(), "RedeemBy() should return correct redeem by time")
	assert.Equal(t, map[string]string{"test": "test"}, coupon.Metadata(), "Metadata() should return correct metadata")
	assert.Equal(t, now, coupon.CreatedAt(), "CreatedAt() should return correct creation time")
	assert.Equal(t, now, coupon.UpdatedAt(), "UpdatedAt() should return correct update time")
	assert.Nil(t, coupon.DeletedAt(), "DeletedAt() should return nil when not set")
}
