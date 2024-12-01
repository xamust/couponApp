package models

import (
	"fmt"
	"github.com/xamust/couponApp/internal/domain/coupon"
	"strings"
	"time"
)

type NewAPICoupon struct {
	Name           string            `json:"name"`
	Reward         string            `json:"reward"`
	MaxRedemptions int               `json:"maxRedemptions"`
	RedeemBy       string            `json:"redeemBy"`
	Metadata       map[string]string `json:"metadata"`
}

func (in *NewAPICoupon) Map() (*coupon.NewCouponInput, error) {
	result := &coupon.NewCouponInput{
		Name:           in.Name,
		Reward:         in.Reward,
		MaxRedemptions: in.MaxRedemptions,
		Metadata:       in.Metadata,
	}
	if strings.TrimSpace(in.RedeemBy) != "" {
		parseTime, err := time.Parse("2006-01-02", in.RedeemBy)
		if err != nil {
			return nil, fmt.Errorf("parse time in api coupon: %w", err)
		}
		result.RedeemBy = &parseTime
	}

	return result, nil
}

type APICoupon struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Reward         string            `json:"reward"`
	MaxRedemptions int               `json:"maxRedemptions"`
	TimesRedeemed  int               `json:"timesRedeemed"`
	RedeemBy       *time.Time        `json:"redeemBy"`
	Metadata       map[string]string `json:"metadata"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
	DeletedAt      *time.Time        `json:"deletedAt"`
}

type APICoupons []APICoupon

func MapCouponResp(usr *coupon.Coupon) *APICoupon {
	return &APICoupon{
		ID:             usr.ID(),
		Name:           usr.Name(),
		Reward:         usr.Reward(),
		MaxRedemptions: usr.MaxRedemptions(),
		TimesRedeemed:  usr.TimesRedeemed(),
		RedeemBy:       usr.RedeemBy(),
		Metadata:       usr.Metadata(),
		CreatedAt:      usr.CreatedAt(),
		UpdatedAt:      usr.UpdatedAt(),
		DeletedAt:      usr.DeletedAt(),
	}
}

func MapCouponRespList(usr []*coupon.Coupon) APICoupons {
	result := make(APICoupons, len(usr))
	for i, v := range usr {
		result[i] = *MapCouponResp(v)
	}
	return result
}

type APICouponList struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type APICouponApplier struct {
	UserID   string `json:"user_id" validate:"required,uuid"`
	CouponID string `json:"coupon_id" validate:"required,uuid"`
}
