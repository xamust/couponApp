package user_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/xamust/couponApp/internal/domain/user"
	"testing"
	"time"
)

func TestNewUserInput_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     *user.NewUserInput
		expectErr error
	}{
		{
			name: "valid input",
			input: &user.NewUserInput{
				Name:     "John Doe",
				Metadata: map[string]string{"age": "30"},
			},
			expectErr: nil,
		},
		{
			name: "empty name",
			input: &user.NewUserInput{
				Name:     " ",
				Metadata: map[string]string{"age": "30"},
			},
			expectErr: user.ErrNameCantBeEmpty,
		},
		{
			name: "nil metadata",
			input: &user.NewUserInput{
				Name:     "John Doe",
				Metadata: nil,
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			assert.Equal(t, tt.expectErr, err)
			if tt.input.Metadata == nil {
				assert.NotNil(t, tt.input.Metadata, "Metadata should be initialized as an empty map")
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	tests := []struct {
		name      string
		input     *user.NewUserInput
		expectErr error
	}{
		{
			name: "create valid user",
			input: &user.NewUserInput{
				Name:     "John Doe",
				Metadata: map[string]string{"age": "30"},
			},
			expectErr: nil,
		},
		{
			name: "invalid input - empty name",
			input: &user.NewUserInput{
				Name:     " ",
				Metadata: map[string]string{"age": "30"},
			},
			expectErr: user.ErrNameCantBeEmpty,
		},
		{
			name: "nil metadata",
			input: &user.NewUserInput{
				Name:     "Jane Doe",
				Metadata: nil,
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := user.New(context.Background(), tt.input)
			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.input.Name, user.Name())
				assert.Equal(t, tt.input.Metadata, user.Metadata())
			}
		})
	}
}

func TestUser_Methods(t *testing.T) {
	now := time.Now()
	usr := &user.RestoreSpecification{
		ID:        "123",
		Name:      "John Doe",
		Metadata:  map[string]string{"age": "30"},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}
	user := usr.RestoreUser()

	assert.Equal(t, "123", user.ID(), "ID() should return correct ID")
	assert.Equal(t, "John Doe", user.Name(), "Name() should return correct name")
	assert.Equal(t, map[string]string{"age": "30"}, user.Metadata(), "Metadata() should return correct metadata")
	assert.Equal(t, now, user.CreatedAt(), "CreatedAt() should return correct creation time")
	assert.Equal(t, now, user.UpdatedAt(), "UpdatedAt() should return correct update time")
	assert.Nil(t, user.DeletedAt(), "DeletedAt() should return nil when not set")
}
