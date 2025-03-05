package models

import (
	"notes/pkg/date"
	"time"
)

// Category represents a note category
// @swagger:model
type Category struct {
	ID        uint       `gorm:"primaryKey" json:"id,omitempty"`
	Name      string     `gorm:"not null;unique;size:30" json:"name"`
	CreatedAt time.Time  `gorm:"created_at" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"updated_at" json:"updated_at,omitempty"`
}

func NewCategory(name string) *Category {
	return &Category{
		Name:      name,
		CreatedAt: *date.ArgentinaTimeNow(),
		UpdatedAt: nil,
	}
}
