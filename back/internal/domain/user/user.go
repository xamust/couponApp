package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

var (
	ErrNameCantBeEmpty = errors.New("name cannot be empty")
	ErrCanNotRemove    = errors.New("can't remove user")
)

type User struct {
	id        string
	name      string
	isActive  bool
	metadata  map[string]string
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func (u *User) ID() string                  { return u.id }
func (u *User) Name() string                { return u.name }
func (u *User) IsActive() bool              { return u.isActive }
func (u *User) Metadata() map[string]string { return u.metadata }
func (u *User) CreatedAt() time.Time        { return u.createdAt }
func (u *User) UpdatedAt() time.Time        { return u.updatedAt }
func (u *User) DeletedAt() *time.Time       { return u.deletedAt }

type NewUserInput struct {
	Name     string
	Metadata map[string]string
}

func (in *NewUserInput) Validate() error {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return ErrNameCantBeEmpty
	}

	if in.Metadata == nil {
		in.Metadata = make(map[string]string)
	}

	return nil
}

func New(_ context.Context, in *NewUserInput) (*User, error) {
	if err := in.Validate(); err != nil {
		return nil, fmt.Errorf("validate new user: %w", err)
	}
	u := &User{
		id:       uuid.New().String(),
		name:     in.Name,
		isActive: true,
		metadata: in.Metadata,
	}
	return u, nil
}

type RestoreSpecification struct {
	ID        string
	Name      string
	IsActive  bool
	Metadata  map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (r *RestoreSpecification) RestoreUser() *User {
	return &User{
		id:        r.ID,
		name:      r.Name,
		isActive:  r.IsActive,
		metadata:  r.Metadata,
		createdAt: r.CreatedAt,
		updatedAt: r.UpdatedAt,
		deletedAt: r.DeletedAt,
	}
}
