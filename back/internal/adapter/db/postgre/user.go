package postgre

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/xamust/couponApp/internal/domain/user"
	"github.com/xamust/couponApp/utils/wrappers"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserRepository struct {
	db_ *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &UserRepository{db_: db}
}

func (rep *UserRepository) db(ctx context.Context) *gorm.DB {
	return rep.db_.WithContext(ctx)
}

func (rep *UserRepository) FindOne(ctx context.Context, id string) (*user.User, error) {
	usr := &User{}
	q := rep.applyCond(rep.db(ctx), user.Cond{Ids: []string{strings.ToLower(id)}})
	if err := q.First(usr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrNotFound
		}
		return nil, err
	}
	return UnmarshalUser(usr)
}

func (rep *UserRepository) Find(ctx context.Context, cond user.Cond, by user.OrderBy, limit, offset int) ([]*user.User, error) {
	models := make([]*User, 0)
	q := rep.applyCond(rep.db(ctx), cond)
	q = rep.applyOrder(q, by)
	q = rep.applyLimit(q, limit, offset)

	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}
	usrs := make([]*user.User, len(models))
	for k, v := range models {
		usr, err := UnmarshalUser(v)
		if err != nil {
			return nil, err
		}
		usrs[k] = usr
	}
	return usrs, nil
}

func (rep *UserRepository) applyCond(q *gorm.DB, cond user.Cond) *gorm.DB {
	q = q.Table("users")

	if len(cond.Ids) > 0 {
		q = q.Where("users.id IN (?)", cond.Ids)
	}
	if len(cond.Metadata) > 0 {
		for k, v := range cond.Metadata {
			q = q.Where("users.metadata->>? IN (?)", k, v)
		}
	}
	return q
}

func (rep *UserRepository) applyOrder(q *gorm.DB, by user.OrderBy) *gorm.DB {
	switch by {
	case user.OrderByCreatedAsc:
		q = q.Order("users.created_at asc")
	default:
		q = q.Order("users.created_at desc")
	}
	return q
}

func (rep *UserRepository) applyLimit(q *gorm.DB, limit, offset int) *gorm.DB {
	if limit > 0 {
		q = q.Limit(limit)
	}
	if limit > 100 {
		q = q.Limit(100)
	}
	return q.Offset(offset)
}

func (rep *UserRepository) Save(ctx context.Context, usr *user.User) error {
	marshalUsr, err := MarshalUser(usr)
	if err != nil {
		return err
	}
	if err := rep.db(ctx).Save(marshalUsr).Error; err != nil {
		return err
	}
	return nil
}

func (rep *UserRepository) Delete(ctx context.Context, usr *user.User) error {
	if err := rep.db_.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		marshalUsr, err := MarshalUser(usr)
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
			return user.ErrCanNotRemove
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

type User struct {
	ID        string `gorm:"PrimaryKey"`
	Name      string
	IsActive  bool
	Metadata  []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (User) TableName() string { return "users" }

func MarshalUser(usr *user.User) (*User, error) {
	md, err := json.Marshal(usr.Metadata())
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        usr.ID(),
		Name:      usr.Name(),
		IsActive:  usr.IsActive(),
		Metadata:  md,
		CreatedAt: usr.CreatedAt(),
		UpdatedAt: usr.UpdatedAt(),
		DeletedAt: wrappers.WrapGormDeletedAt(usr.DeletedAt()),
	}, nil
}

func UnmarshalUser(usr *User) (*user.User, error) {
	md := map[string]string{}
	if usr.Metadata != nil {
		if err := json.Unmarshal(usr.Metadata, &md); err != nil {
			return nil, err
		}
	}
	restoreUser := user.RestoreSpecification{
		ID:        usr.ID,
		Name:      usr.Name,
		IsActive:  usr.IsActive,
		Metadata:  md,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		DeletedAt: wrappers.UnwrapGormDeletedAt(usr.DeletedAt),
	}
	return restoreUser.RestoreUser(), nil
}
