package models

import (
	"github.com/xamust/couponApp/internal/domain/user"
	"time"
)

type NewAPIUser struct {
	Name     string            `json:"name" example:"John Doe"`
	Metadata map[string]string `json:"metadata" example:"{\"key\":\"value\"}"`
}

func (in *NewAPIUser) Map() *user.NewUserInput {
	return &user.NewUserInput{
		Name:     in.Name,
		Metadata: in.Metadata,
	}
}

type APIUser struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	IsActive  bool              `json:"is_active"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt *time.Time        `json:"deleted_at"`
}

type APIUsers []APIUser

func MapUserResp(usr *user.User) *APIUser {
	return &APIUser{
		ID:        usr.ID(),
		Name:      usr.Name(),
		IsActive:  usr.IsActive(),
		Metadata:  usr.Metadata(),
		CreatedAt: usr.CreatedAt(),
		UpdatedAt: usr.UpdatedAt(),
		DeletedAt: usr.DeletedAt(),
	}
}

func MapUserRespList(usr []*user.User) APIUsers {
	result := make(APIUsers, len(usr))
	for i, v := range usr {
		result[i] = *MapUserResp(v)
	}
	return result
}

type APIUserList struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
