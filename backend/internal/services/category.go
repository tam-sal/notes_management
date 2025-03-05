package services

import (
	"context"
	"notes/internal/models"
	"notes/internal/repositories"
	"notes/pkg/date"
	"notes/pkg/utils"
	"notes/pkg/validations"

	"gorm.io/gorm"
)

type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: repo,
	}
}

func (cs *CategoryService) Create(ctx context.Context, name string) (*models.Category, error) {
	valid, formattedName, err := utils.ValidateAndFormatCategory(name)
	if !valid {
		return nil, err
	}

	c, err := cs.categoryRepo.FindByName(ctx, formattedName)
	if c != nil {
		return nil, validations.ErrCatAlreadyExist
	}
	if err != gorm.ErrRecordNotFound {
		return nil, validations.ErrFetchingCategory
	}
	category := models.NewCategory(formattedName)
	createdCategory, err := cs.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, validations.ErrCatCreate
	}
	return createdCategory, nil
}

func (cs *CategoryService) Update(ctx context.Context, id uint, newName string) (*models.Category, error) {
	valid, formattedName, err := utils.ValidateAndFormatCategory(newName)
	if !valid {
		return nil, err
	}

	category, err := cs.categoryRepo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, validations.ErrCatNotFound
		}
		return nil, validations.ErrFetchingCategory
	}

	existingCategory, err := cs.categoryRepo.FindByName(ctx, formattedName)
	if err == nil && existingCategory.ID != id {
		return nil, validations.ErrCatAlreadyExist
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, validations.ErrFetchingCategory
	}

	category.Name = formattedName
	category.UpdatedAt = date.ArgentinaTimeNow()

	updatedID, err := cs.categoryRepo.Update(ctx, category)
	if err != nil {
		return nil, validations.ErrCatUpdate
	}
	category.ID = *updatedID
	return category, nil
}

func (cs *CategoryService) Delete(ctx context.Context, id uint) (*uint, error) {
	_, err := cs.categoryRepo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, validations.ErrCatNotFound
		}
		return nil, validations.ErrFetchingCategory
	}

	deletedID, err := cs.categoryRepo.Delete(ctx, id)
	if err != nil {
		return nil, validations.ErrCatDelete
	}
	return deletedID, nil
}

func (cs *CategoryService) GetAll(ctx context.Context) ([]models.Category, error) {
	categories, err := cs.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, validations.ErrFetchingCategories
	}
	return categories, nil
}

func (cs *CategoryService) GetById(ctx context.Context, id uint) (*models.Category, error) {
	category, err := cs.categoryRepo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, validations.ErrCatNotFound
		}
		return nil, validations.ErrFetchingCategory
	}
	return category, nil
}

func (cs *CategoryService) GetByName(ctx context.Context, name string) (*models.Category, error) {
	valid, formattedName, err := utils.ValidateAndFormatCategory(name)
	if !valid {
		return nil, err
	}

	category, err := cs.categoryRepo.FindByName(ctx, formattedName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, validations.ErrCatNotFound
		}
		return nil, validations.ErrFetchingCategory
	}
	return category, nil
}

func (cs *CategoryService) GetByNameOrCreate(ctx context.Context, name string) (*models.Category, error) {
	valid, formattedName, err := utils.ValidateAndFormatCategory(name)
	if !valid {
		return nil, err
	}

	category, err := cs.categoryRepo.FindByName(ctx, formattedName)
	if category != nil {
		return category, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, validations.ErrFetchingCategory
	}

	category = models.NewCategory(formattedName)
	createdCategory, err := cs.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, validations.ErrCatCreate
	}
	return createdCategory, nil
}
