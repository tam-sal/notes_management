package repositories

import (
	"context"
	"notes/internal/configs"
	"notes/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db     *gorm.DB
	config *configs.Config
}

func NewCategoryRepository(db *gorm.DB, config *configs.Config) *CategoryRepository {
	return &CategoryRepository{
		db:     db,
		config: config,
	}
}

func (cr *CategoryRepository) Create(ctx context.Context, category *models.Category) (*models.Category, error) {
	if err := cr.db.WithContext(ctx).Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (cr *CategoryRepository) Update(ctx context.Context, category *models.Category) (*uint, error) {
	if err := cr.db.WithContext(ctx).Save(category).Error; err != nil {
		return nil, err
	}
	return &category.ID, nil
}

func (cr *CategoryRepository) Delete(ctx context.Context, id uint) (*uint, error) {
	if err := cr.db.WithContext(ctx).Delete(&models.Category{}, id).Error; err != nil {
		return nil, err
	}
	return &id, nil
}

func (cr *CategoryRepository) FindByID(ctx context.Context, id uint) (*models.Category, error) {
	var category models.Category
	if err := cr.db.WithContext(ctx).First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (cr *CategoryRepository) FindByName(ctx context.Context, name string) (*models.Category, error) {
	var category models.Category
	if err := cr.db.WithContext(ctx).Where("name = ?", name).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (cr *CategoryRepository) FindAll(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	if err := cr.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
