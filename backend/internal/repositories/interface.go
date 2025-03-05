package repositories

import (
	"context"
	"notes/internal/models"
)

type UserRepositoryMethods interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}
