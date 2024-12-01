package user

import (
	"context"
	"errors"
)

type OrderBy string

const (
	OrderByCreatedAsc  OrderBy = "created.asc"
	OrderByCreatedDesc OrderBy = "created.desc"
)

var ErrNotFound = errors.New("user not found")

type Repository interface {
	FindOne(ctx context.Context, id string) (*User, error)
	Find(ctx context.Context, cond Cond, by OrderBy, limit, offset int) ([]*User, error)
	Save(ctx context.Context, co *User) error
	Delete(ctx context.Context, coup *User) error
}

type Cond struct {
	Ids      []string
	Metadata map[string]interface{}
}
