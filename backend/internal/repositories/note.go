package repositories

import (
	"context"
	"fmt"
	"notes/internal/configs"
	"notes/internal/models"
	"notes/pkg/validations"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type NoteRepository struct {
	db           *gorm.DB
	config       *configs.Config
	categoryRepo *CategoryRepository
}

func NewNoteRepository(db *gorm.DB, config *configs.Config, categoryRepo *CategoryRepository) *NoteRepository {
	return &NoteRepository{
		db:           db,
		config:       config,
		categoryRepo: categoryRepo,
	}
}

func (nr *NoteRepository) Create(ctx context.Context, note *models.Note) (*models.Note, error) {
	if note.UserID == 0 {
		return nil, validations.ErrUserIdNotSet
	}
	if err := nr.db.WithContext(ctx).Create(&note).Error; err != nil {
		return nil, validations.ErrNoteCreate
	}
	return note, nil
}

func (nr *NoteRepository) GetAllNotes(ctx context.Context) ([]*models.Note, error) {
	var notes []*models.Note

	if err := nr.db.WithContext(ctx).Preload("Categories").Find(&notes).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, validations.ErrNoNotesFound
		}
		return nil, validations.ErrFetchingNotes
	}

	return notes, nil
}

func (nr *NoteRepository) GetNotesByCategories(ctx context.Context, categoryNames []string) ([]*models.Note, error) {
	var notes []*models.Note
	if err := nr.db.WithContext(ctx).
		Joins("JOIN note_categories ON notes.id = note_categories.note_id").
		Joins("JOIN categories ON note_categories.category_id = categories.id").
		Where("categories.name IN ?", categoryNames).
		Preload("Categories").
		Group("notes.id").
		Find(&notes).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, validations.ErrNoNotesFound
		}
		return nil, validations.ErrFetchingNotes
	}

	return notes, nil
}

func (nr *NoteRepository) GetNoteById(ctx context.Context, id uint) (*models.Note, error) {

	var note models.Note

	if err := nr.db.WithContext(ctx).Preload("Categories").First(&note, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, validations.ErrNoteNotFound
		}
		return nil, validations.ErrFetchingNote
	}
	return &note, nil
}

func (nr *NoteRepository) GetByTitle(ctx context.Context, title string) (*models.Note, error) {
	var note models.Note
	if err := nr.db.WithContext(ctx).Where("title = ?", title).First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, validations.ErrNotTitle
		}
		return nil, err
	}

	return &note, nil
}

func (nr *NoteRepository) UpdateNote(ctx context.Context, note *models.Note) (*uint, error) {
	tx := nr.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := tx.Model(note).Association("Categories").Replace(note.Categories); err != nil {
		return nil, validations.ErrCatUpdate
	}

	if err := tx.Save(note).Error; err != nil {
		return nil, validations.ErrNoteUpdate
	}

	if err := tx.Commit().Error; err != nil {
		return nil, validations.ErrNoteUpdate
	}

	return &note.ID, nil
}

func (nr *NoteRepository) Delete(ctx context.Context, note *models.Note) (*uint, error) {
	noteId := note.ID
	if err := nr.db.WithContext(ctx).Select("Categories").Delete(&note).Error; err != nil {
		return nil, validations.ErrNoteDelete
	}
	return &noteId, nil
}

func (nr *NoteRepository) FilterNotes(ctx context.Context, isArchived *bool, categories []string) ([]models.Note, error) {
	var notes []models.Note
	query := nr.db.WithContext(ctx).Model(&models.Note{})

	if isArchived != nil {
		query = query.Where("is_archived = ?", *isArchived)
	}

	if len(categories) > 0 {
		query = query.Joins("JOIN note_categories ON notes.id = note_categories.note_id").
			Joins("JOIN categories ON categories.id = note_categories.category_id").
			Where("categories.name ILIKE ANY (CAST(? AS varchar[]))", pq.Array(categories)).
			Group("notes.id")
	}

	if err := query.Preload("Categories").Find(&notes).Error; err != nil {
		fmt.Println("Error during filtering:", err)
		return nil, validations.ErrFilterDB
	}

	if len(notes) == 0 {
		return nil, validations.ErrNoNotesFound
	}

	return notes, nil
}

func (nr *NoteRepository) DeleteNoteByUserId(ctx context.Context, noteId uint, userId uint) (*uint, error) {
	var note models.Note
	if err := nr.db.WithContext(ctx).Where("id = ? AND user_id = ?", noteId, userId).First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, validations.ErrNoteNotOwnedByUser
		}
		return nil, validations.ErrNoteDelete
	}
	if _, err := nr.Delete(ctx, &note); err != nil {
		return nil, err
	}
	return &note.ID, nil
}
