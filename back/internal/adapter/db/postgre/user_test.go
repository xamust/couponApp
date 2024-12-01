package postgre_test

import (
	"context"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/xamust/couponApp/internal/adapter/db/postgre"
	dbMOCK "github.com/xamust/couponApp/internal/adapter/db/postgre/mock"
	"github.com/xamust/couponApp/internal/domain/user"
	"testing"
)

func TestUserRepository_FindOne(t *testing.T) {
	db, mock, err := dbMOCK.SetupMockDB()
	assert.NoError(t, err)

	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	ctx := context.Background()
	repo := postgre.NewUserRepository(db)

	expectedID := "test-id"
	mockCoupon := &postgre.User{
		ID:   expectedID,
		Name: "Test User",
	}

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE users\.id IN \(\$1\) AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(expectedID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(mockCoupon.ID, mockCoupon.Name))

	result, err := repo.FindOne(ctx, expectedID)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, result.ID())
	assert.Equal(t, mockCoupon.Name, result.Name())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Find(t *testing.T) {
	db, mock, err := dbMOCK.SetupMockDB()
	assert.NoError(t, err)

	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	ctx := context.Background()
	repo := postgre.NewUserRepository(db)

	cond := user.Cond{
		Ids: []string{"test-id"},
	}

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE users\.id IN \(\$1\) AND "users"\."deleted_at" IS NULL ORDER BY users\.created_at asc LIMIT \$2`).
		WithArgs("test-id", 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow("test-id", "Test Coupon"))

	results, err := repo.Find(ctx, cond, user.OrderByCreatedAsc, 10, 0)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "test-id", results[0].ID())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Save(t *testing.T) {
	db, mock, err := dbMOCK.SetupMockDB()
	assert.NoError(t, err)

	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	ctx := context.Background()
	repo := postgre.NewUserRepository(db)

	md := map[string]interface{}{"someKey": "someValue"}
	marshalMD, err := json.Marshal(md)
	assert.NoError(t, err)

	mockUser := &postgre.User{
		ID:       "new-id",
		Name:     "New Coupon",
		IsActive: true,
		Metadata: marshalMD,
	}

	coup, err := postgre.UnmarshalUser(mockUser)
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users"`). //`INSERT INTO "users"`
						WithArgs(
			mockUser.Name,
			mockUser.IsActive,
			mockUser.Metadata,
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt (может быть nil)
			mockUser.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Save(ctx, coup)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, err := dbMOCK.SetupMockDB()
	assert.NoError(t, err)

	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	ctx := context.Background()
	repo := postgre.NewUserRepository(db)

	userToDelete := &postgre.User{
		ID: "delete-id",
	}

	expectedID := "delete-id"

	coup, err := postgre.UnmarshalUser(userToDelete)
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users" SET "deleted_at"=\$1,"updated_at"=\$2 WHERE "users"\."deleted_at" IS NULL AND "id" = \$3`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), expectedID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Delete(ctx, coup)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
