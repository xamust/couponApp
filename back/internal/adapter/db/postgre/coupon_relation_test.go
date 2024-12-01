package postgre_test

import (
	"context"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/xamust/couponApp/internal/adapter/db/postgre"
	dbMOCK "github.com/xamust/couponApp/internal/adapter/db/postgre/mock"
	"github.com/xamust/couponApp/internal/domain/coupon_relation"
	"testing"
)

func TestCouponRelationRepository_Create(t *testing.T) {
	db, mock, err := dbMOCK.SetupMockDB()
	assert.NoError(t, err)

	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	ctx := context.Background()
	repo := postgre.NewCouponRelationRepository(db)

	md := map[string]interface{}{"someKey": "someValue"}
	marshalMD, err := json.Marshal(md)
	assert.NoError(t, err)

	mockCouponApplier := &postgre.CouponRelation{
		ID:       "new-id",
		UserID:   "user-id",
		CouponID: "coupon-id",
		Metadata: marshalMD,
	}

	coup, err := postgre.UnmarshalCouponRelation(mockCouponApplier)
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "coupon_relation" \("id","user_id","coupon_id","metadata","created_at","updated_at","deleted_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7\)`). //`INSERT INTO "coupon_relation"`
																							WithArgs(
			mockCouponApplier.ID,
			mockCouponApplier.UserID,
			mockCouponApplier.CouponID,
			mockCouponApplier.Metadata,
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt (может быть nil)

		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Create(ctx, coup)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCouponRelationRepository_FindOne(t *testing.T) {
	db, mock, err := dbMOCK.SetupMockDB()
	assert.NoError(t, err)

	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	ctx := context.Background()
	repo := postgre.NewCouponRelationRepository(db)

	expectedID := "test-id"
	mockCouponRelation := &postgre.CouponRelation{
		ID:       expectedID,
		UserID:   "user-id",
		CouponID: "coupon-id",
	}

	mock.ExpectQuery(`SELECT \* FROM "coupon_relation" WHERE coupon_relation\.id IN \(\$1\) AND "coupon_relation"\."deleted_at" IS NULL ORDER BY "coupon_relation"\."id" LIMIT \$2`).
		WithArgs(expectedID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "coupon_id"}).
			AddRow(
				mockCouponRelation.ID,
				mockCouponRelation.UserID,
				mockCouponRelation.CouponID))

	result, err := repo.FindOne(ctx, expectedID)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, result.ID())
	assert.Equal(t, mockCouponRelation.UserID, result.UserID())
	assert.Equal(t, mockCouponRelation.CouponID, result.CouponID())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCouponRelationRepository_Find(t *testing.T) {
	db, mock, err := dbMOCK.SetupMockDB()
	assert.NoError(t, err)

	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	ctx := context.Background()
	repo := postgre.NewCouponRelationRepository(db)

	cond := coupon_relation.Cond{
		Ids: []string{"test-id-1", "test-id-2"},
	}

	mock.ExpectQuery(`SELECT \* FROM "coupon_relation" WHERE coupon_relation\.id IN \(\$1,\$2\) AND "coupon_relation"\."deleted_at" IS NULL ORDER BY coupon_relation\.created_at asc LIMIT \$3`).
		WithArgs("test-id-1", "test-id-2", 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "coupon_id"}).
			AddRow("test-id-1", "user-id-1", "coupon-id-1").
			AddRow("test-id-2", "user-id-2", "coupon-id-2"))
	//todo @ check on real DB
	results, err := repo.Find(ctx, cond, coupon_relation.OrderByCreatedAsc, 10, 0)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "test-id-1", results[0].ID())
	assert.Equal(t, "test-id-2", results[1].ID())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCouponRelationRepository_Delete(t *testing.T) {
	db, mock, err := dbMOCK.SetupMockDB()
	assert.NoError(t, err)

	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	ctx := context.Background()
	repo := postgre.NewCouponRelationRepository(db)

	couponRelationToDelete := &postgre.CouponRelation{
		ID: "delete-id",
	}

	expectedID := "delete-id"

	coup, err := postgre.UnmarshalCouponRelation(couponRelationToDelete)
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "coupon_relation" SET "deleted_at"=\$1,"updated_at"=\$2 WHERE "coupon_relation"\."deleted_at" IS NULL AND "id" = \$3`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), expectedID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Delete(ctx, coup)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
