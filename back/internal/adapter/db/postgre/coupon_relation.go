package postgre

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xamust/couponApp/internal/domain/coupon_relation"
	"github.com/xamust/couponApp/internal/domain/user"
	"github.com/xamust/couponApp/utils/wrappers"
	"gorm.io/gorm"
	"strings"
	"time"
)

type CouponRelationRepository struct {
	db_ *gorm.DB
}

func NewCouponRelationRepository(db *gorm.DB) coupon_relation.Repository {
	return &CouponRelationRepository{db_: db}
}

func (rep *CouponRelationRepository) db(ctx context.Context) *gorm.DB {
	return rep.db_.WithContext(ctx)
}

func (rep *CouponRelationRepository) Create(ctx context.Context, ca *coupon_relation.CouponRelation) error {
	marshalCA, err := MarshalCouponRelation(ca)
	if err != nil {
		return err
	}
	if err := rep.db(ctx).Create(marshalCA).Error; err != nil {
		return err
	}
	return nil
}

func (rep *CouponRelationRepository) FindOne(ctx context.Context, ID string) (*coupon_relation.CouponRelation, error) {
	ca := &CouponRelation{}
	q := rep.applyCond(rep.db(ctx), coupon_relation.Cond{Ids: []string{strings.ToLower(ID)}})
	if err := q.First(ca).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrNotFound
		}
		return nil, err
	}
	return UnmarshalCouponRelation(ca)
}

func (rep *CouponRelationRepository) Find(ctx context.Context, cond coupon_relation.Cond, by coupon_relation.OrderBy, limit, offset int) ([]*coupon_relation.CouponRelation, error) {
	models := make([]CouponRelation, 0)
	q := rep.applyCond(rep.db(ctx), cond)
	q = rep.applyOrder(q, by)
	q = rep.applyLimit(q, limit, offset)
	fmt.Println(q.Statement.SQL.String())
	if err := q.Debug().Find(&models).Error; err != nil { // todo @ debug to off
		return nil, err
	}
	ca := make([]*coupon_relation.CouponRelation, len(models))
	for k, v := range models {
		usr, err := UnmarshalCouponRelation(&v)
		if err != nil {
			return nil, err
		}
		ca[k] = usr
	}
	return ca, nil
}

func (rep *CouponRelationRepository) applyCond(q *gorm.DB, cond coupon_relation.Cond) *gorm.DB {
	q = q.Table("coupon_relation")

	if len(cond.Ids) > 0 {
		q = q.Where("coupon_relation.id IN (?)", cond.Ids)
	}
	if len(cond.UserIDs) > 0 {
		q = q.Where("coupon_relation.user_id IN (?)", cond.UserIDs)
	}
	if len(cond.CouponIDs) > 0 {
		q = q.Where("coupon_relation.coupon_id IN (?)", cond.CouponIDs)
	}
	if len(cond.Metadata) > 0 {
		for k, v := range cond.Metadata {
			q = q.Where("coupon_relation.metadata->>? IN (?)", k, v)
		}
	}
	return q
}

func (rep *CouponRelationRepository) applyOrder(q *gorm.DB, by coupon_relation.OrderBy) *gorm.DB {
	switch by {
	case coupon_relation.OrderByCreatedAsc:
		q = q.Order("coupon_relation.created_at asc")
	default:
		q = q.Order("coupon_relation.created_at desc")
	}
	return q
}

func (rep *CouponRelationRepository) applyLimit(q *gorm.DB, limit, offset int) *gorm.DB {
	if limit > 0 {
		q = q.Limit(limit)
	}
	if limit > 100 {
		q = q.Limit(100)
	}
	return q.Offset(offset)
}

func (rep *CouponRelationRepository) Delete(ctx context.Context, ca *coupon_relation.CouponRelation) error {
	if err := rep.db_.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		marshalUsr, err := MarshalCouponRelation(ca)
		if err != nil {
			tx.Rollback()
			return err
		}
		callback := tx.Model(marshalUsr).Update("deleted_at", time.Now())
		if callback.Error != nil {
			return callback.Error
		}
		if callback.RowsAffected == 0 {
			tx.Rollback()
			return coupon_relation.ErrCanNotRemove
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

type CouponRelation struct {
	ID        string `gorm:"PrimaryKey"`
	UserID    string
	CouponID  string
	Metadata  []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (CouponRelation) TableName() string { return "coupon_relation" }

func MarshalCouponRelation(ca *coupon_relation.CouponRelation) (*CouponRelation, error) {
	md, err := json.Marshal(ca.Metadata())
	if err != nil {
		return nil, err
	}
	return &CouponRelation{
		ID:        ca.ID(),
		UserID:    ca.UserID(),
		CouponID:  ca.CouponID(),
		Metadata:  md,
		CreatedAt: ca.CreatedAt(),
		UpdatedAt: ca.UpdatedAt(),
		DeletedAt: wrappers.WrapGormDeletedAt(ca.DeletedAt()),
	}, nil
}

func UnmarshalCouponRelation(ca *CouponRelation) (*coupon_relation.CouponRelation, error) {
	md := map[string]interface{}{}
	if ca.Metadata != nil {
		if err := json.Unmarshal(ca.Metadata, &md); err != nil {
			return nil, err
		}
	}
	restoreCA := coupon_relation.RestoreSpecification{
		ID:        ca.ID,
		UserID:    ca.UserID,
		CouponID:  ca.CouponID,
		Metadata:  md,
		CreatedAt: ca.CreatedAt,
		UpdatedAt: ca.UpdatedAt,
		DeletedAt: wrappers.UnwrapGormDeletedAt(ca.DeletedAt),
	}
	return restoreCA.RestoreCouponRelation(), nil
}
