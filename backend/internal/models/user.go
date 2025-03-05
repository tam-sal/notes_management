package models

import (
	"notes/pkg/date"
	"time"
)

// User represents a user
// @swagger:model
type User struct {
	ID        uint       `gorm:"primaryKey" json:"id,omitempty"`
	UserName  string     `gorm:"unique;not null;size:30" json:"user_name"`
	Password  string     `gorm:"not null" json:"-"`
	Notes     []Note     `gorm:"constraint:OnDelete:CASCADE;" json:"notes"`
	IsAdmin   bool       `gorm:"default:false" json:"is_admin,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

func NewUser(user, password string, notes []*Note) *User {

	return &User{
		UserName:  user,
		Password:  password,
		Notes:     nil,
		IsAdmin:   false,
		CreatedAt: *date.ArgentinaTimeNow(),
		UpdatedAt: nil,
	}
}
