package postgre_test

import (
	"context"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/xamust/couponApp/internal/adapter/db/postgre"
	dbMOCK "github.com/xamust/couponApp/internal/adapter/db/postgre/mock"
	"github.com/xamust/couponApp/internal/domain/coupon"
	"regexp"
	"testing"
)

func TestCouponRepository_FindOne(t *testing.T) {
	db, mock, err := dbMOCK.SetupMockDB()
	assert.NoError(t, err)

	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	ctx := context.Background()
	repo := postgre.NewCouponRepository(db)

	expectedID := "test-id"
	mockCoupon := &postgre.Coupon{
		ID:             expectedID,
		Name:           "Test Coupon",
		Reward:         "Free Coffee",
		MaxRedemptions: 10,
		TimesRedeemed:  1,
		RedeemBy:       nil,
	}

	mock.ExpectQuery(`SELECT \* FROM "coupons" WHERE coupons\.id IN \(\$1\) AND "coupons"\."deleted_at" IS NULL ORDER BY "coupons"\."id" LIMIT \$2`).
		WithArgs(expectedID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "reward", "max_redemptions", "times_redeemed"}).
			AddRow(mockCoupon.ID, mockCoupon.Name, mockCoupon.Reward, mockCoupon.MaxRedemptions, mockCoupon.TimesRedeemed))

	result, err := repo.FindOne(ctx, expectedID)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, result.ID())
	assert.Equal(t, mockCoupon.Name, result.Name())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCouponRepository_Delete(t *testing.T) {
	db, mock, err := dbMOCK.SetupMockDB()
	assert.NoError(t, err)

	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	ctx := context.Background()
	repo := postgre.NewCouponRepository(db)

	couponToDelete := &postgre.Coupon{
		ID: "delete-id",
	}

	expectedID := "delete-id"

	coup, err := postgre.UnmarshalCoupon(couponToDelete)
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "coupons" SET "deleted_at"=\$1,"updated_at"=\$2 WHERE times_redeemed = 0 AND "coupons"\."deleted_at" IS NULL AND "id" = \$3`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), expectedID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Delete(ctx, coup)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCouponRepository_Save(t *testing.T) {
	db, mock, err := dbMOCK.SetupMockDB()
	assert.NoError(t, err)

	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	ctx := context.Background()
	repo := postgre.NewCouponRepository(db)

	md := map[string]string{"someKey": "someValue"}
	marshalMD, err := json.Marshal(md)
	assert.NoError(t, err)

	mockCoupon := &postgre.Coupon{
		Name:           "New Coupon",
		Reward:         "50% Discount",
		MaxRedemptions: 100,
		TimesRedeemed:  0,
		Metadata:       marshalMD,
	}

	coup, err := postgre.UnmarshalCoupon(mockCoupon)
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec(
		regexp.QuoteMeta(`INSERT INTO "coupons" ("id","name","reward","max_redemptions","times_redeemed","redeem_by","metadata","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`),
	).
		WithArgs(
			"",
			mockCoupon.Name,
			mockCoupon.Reward,
			mockCoupon.MaxRedemptions,
			mockCoupon.TimesRedeemed,
			sqlmock.AnyArg(), // RedeemBy может быть nil
			mockCoupon.Metadata,
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt (может быть nil)
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = repo.Save(ctx, coup)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCouponRepository_Find(t *testing.T) {
	db, mock, err := dbMOCK.SetupMockDB()
	assert.NoError(t, err)

	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	ctx := context.Background()
	repo := postgre.NewCouponRepository(db)

	cond := coupon.Cond{
		Ids: []string{"test-id"},
	}

	mock.ExpectQuery(`SELECT \* FROM "coupons" WHERE coupons\.id IN \(\$1\) AND "coupons"\."deleted_at" IS NULL ORDER BY coupons\.created_at asc LIMIT \$2`).
		WithArgs("test-id", 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "reward"}).
			AddRow("test-id", "Test Coupon", "50% Discount"))

	results, err := repo.Find(ctx, cond, coupon.OrderByCreatedAsc, 10, 0)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "test-id", results[0].ID())

	assert.NoError(t, mock.ExpectationsWereMet())
}
