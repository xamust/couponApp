package postgre_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xamust/couponApp/internal/adapter/db/postgre"
	dbMOCK "github.com/xamust/couponApp/internal/adapter/db/postgre/mock"
	"github.com/xamust/couponApp/internal/domain/coupon"
	"github.com/xamust/couponApp/internal/domain/coupon_relation"
	"testing"
)

func TestApply(t *testing.T) {
	ctx := context.Background()
	db, mock, err := dbMOCK.SetupMockDB()
	assert.NoError(t, err)

	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	repo := postgre.NewCouponApplierRepository(db)

	mockCoupon := &postgre.Coupon{
		ID:             "test-coupon-id",
		Name:           "test-coupon-name",
		Reward:         "test-coupon-reward",
		MaxRedemptions: 0,
		TimesRedeemed:  0,
	}

	expectedID := "test-relation-id"
	couponRelationMock := &postgre.CouponRelation{
		ID:       expectedID,
		UserID:   "test-user-id",
		CouponID: "test-coupon-id",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "coupons"`).
		WithArgs(
			mockCoupon.Name,
			mockCoupon.Reward,
			mockCoupon.MaxRedemptions,
			mockCoupon.TimesRedeemed,
			sqlmock.AnyArg(), // RedeemBy может быть nil
			sqlmock.AnyArg(),
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt (может быть nil)
			mockCoupon.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO "coupon_relation" \("id","user_id","coupon_id","metadata","created_at","updated_at","deleted_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7\)`).
		WithArgs(
			couponRelationMock.ID,
			couponRelationMock.UserID,
			couponRelationMock.CouponID,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	c := coupon.RestoreSpecification{
		ID:             "test-coupon-id",
		Name:           "test-coupon-name",
		Reward:         "test-coupon-reward",
		MaxRedemptions: 0,
		TimesRedeemed:  0,
	}
	coup := c.RestoreCoupon()

	r := coupon_relation.RestoreSpecification{
		ID:       "test-relation-id",
		UserID:   "test-user-id",
		CouponID: "test-coupon-id",
	}
	cr := r.RestoreCouponRelation()

	err = repo.CouponApplier(ctx, coup, cr)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}
