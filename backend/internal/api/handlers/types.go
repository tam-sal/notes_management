package handlers

import (
	"notes/internal/models"
	"time"
)

// CreateNoteRequest represents the payload for creating a new note
// @swagger:model CreateNoteRequest
type CreateNoteRequest struct {
	Title      string   `json:"title" example:"My Note Title"`
	Content    string   `json:"content" example:"This is the content of the note."`
	Categories []string `json:"categories" example:"[\"Work\", \"Personal\"]"`
	UserID     uint     `json:"user_id" example:"1"`
}

// GetNoteResponse represents the response when fetching a note
// @swagger:model GetNoteResponse
type GetNoteResponse struct {
	ID         uint              `json:"id" example:"1"`
	Title      string            `json:"title" example:"Sample Note Title"`
	Content    string            `json:"content" example:"Sample note content."`
	Categories []models.Category `json:"categories"`
	IsArchived bool              `json:"is_archived" example:"false"`
	CreatedAt  time.Time         `json:"created_at" example:"2025-02-01T12:00:00Z"`
	UpdatedAt  *time.Time        `json:"updated_at,omitempty"`
	UserID     uint              `json:"user_id" example:"1"`
}

// UpdateNoteRequest represents the payload for updating a note
// @swagger:model
type UpdateNoteRequest struct {
	Title      string            `json:"title"`
	Content    string            `json:"content"`
	Categories []models.Category `json:"categories"`
	IsArchived bool              `json:"is_archived"`
	UserID     uint              `json:"user_id"`
}

// CreateCategoryRequest represents the payload for creating a new category.
// @Description Payload for creating a new category
type CreateCategoryRequest struct {
	Name string `json:"name"`
}

// GetCategory represents the response when fetching a category.
// @Description Response containing details of a category
type GetCategory struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// UserRequest represents registration request
// @swagger:model
type UserRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// StatusResponse represents the response for the /status endpoint.
type StatusResponse struct {
	Status string `json:"status" example:"OK"`
}

// APIResponse represents the standard structure of API responses
// @swagger:model APIResponse
type APIResponse struct {
	Error  string      `json:"error,omitempty" example:"Error message if any"`
	Data   interface{} `json:"data,omitempty"`
	Status int         `json:"status" example:"200"`
}

// UserResponse represents successful registration response
// @swagger:model
type UserResponse struct {
	ID       *uint  `json:"id"`
	UserName string `json:"user_name"`

	// Notes associated with the user
	// swagger:allOf []models.Note
	Notes []models.Note `json:"notes"`

	IsAdmin   bool       `json:"is_admin"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// ErrorResponse represents error format
// @swagger:model
type ErrorResponse struct {
	// Error message
	// Example: Invalid input data
	Error string `json:"error"`

	// No data - null
	// Example: nil
	Data interface{} `json:"data,omitempty"`

	// HTTP status code
	// Example: 400
	Status int `json:"status"`
}

// @swagger:model
type LogOut struct {
	Data   interface{} `json:"data"`
	Error  string      `json:"error,omitempty"`
	Status int         `json:"status"`
}
