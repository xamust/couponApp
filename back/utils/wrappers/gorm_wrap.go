package wrappers

import (
	"gorm.io/gorm"
	"time"
)

func WrapGormDeletedAt(t *time.Time) gorm.DeletedAt {
	deleted := gorm.DeletedAt{}
	if t != nil {
		deleted = gorm.DeletedAt{
			Time:  *t,
			Valid: true,
		}
	}
	return deleted
}

func UnwrapGormDeletedAt(t gorm.DeletedAt) *time.Time {
	var deleted *time.Time
	if t.Valid {
		deleted = &t.Time
	}
	return deleted
}
