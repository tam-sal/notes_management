package repositories

import (
	"context"
	"notes/internal/configs"
	"notes/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db     *gorm.DB
	config *configs.Config
}

func NewUserRepository(db *gorm.DB, config *configs.Config) *UserRepository {
	return &UserRepository{
		db:     db,
		config: config,
	}
}
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.Create(user).Error
}

func (ur *UserRepository) GetUserByID(ctx context.Context, userId uint) (*models.User, error) {
	var user models.User
	if err := ur.db.WithContext(ctx).
		Preload("Notes.Categories").
		First(&user, userId).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := ur.db.WithContext(ctx).Preload("Notes.Categories").Where("user_name ILIKE ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
