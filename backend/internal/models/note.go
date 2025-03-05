package models

import (
	"notes/pkg/date"
	"time"
)

// Note represents a user's note
// @swagger:model
type Note struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Title      string     `gorm:"not null;unique;size:50" json:"title"`
	Content    string     `gorm:"type:text;size:70" json:"content"`
	Categories []Category `gorm:"many2many:note_categories;" json:"categories"`
	UserID     uint       `gorm:"not null;foreignKey:UserID" json:"user_id"`
	IsArchived bool       `gorm:"default:false" json:"is_archived"`
	CreatedAt  time.Time  `gorm:"created_at" json:"created_at,omitempty"`
	UpdatedAt  *time.Time `gorm:"updated_at" json:"updated_at,omitempty"`
}

func NewNote(title, content string, categories []Category, userId uint) *Note {
	return &Note{
		Title:      title,
		Content:    content,
		Categories: categories,
		UserID:     userId,
		IsArchived: false,
		CreatedAt:  *date.ArgentinaTimeNow(),
		UpdatedAt:  nil,
	}
}
