package services

import (
	"context"
	"notes/internal/models"
	"notes/internal/repositories"
	"notes/pkg/date"
	"notes/pkg/utils"
	"notes/pkg/validations"
)

type NoteService struct {
	noteRepo        *repositories.NoteRepository
	CategoryService *CategoryService
}

func NewNoteService(noteRepo *repositories.NoteRepository, categoryService *CategoryService) *NoteService {
	return &NoteService{
		noteRepo:        noteRepo,
		CategoryService: categoryService,
	}
}

func (ns *NoteService) CreateNote(ctx context.Context, title, content string, categoryNames []string, userID uint) (*models.Note, error) {
	valid, formattedTitle, err := utils.ValidateAndFormatTitle(title)
	if !valid {
		return nil, err
	}

	valid, formattedContent, err := utils.ValidateAndFormatContent(content)
	if !valid {
		return nil, err
	}

	if len(categoryNames) > 4 {
		return nil, validations.ErrTooManyCategories
	}
	if len(categoryNames) < 1 {
		return nil, validations.ErrZeroCategory
	}

	var categories []models.Category
	for _, categoryName := range categoryNames {
		category, err := ns.CategoryService.GetByNameOrCreate(ctx, categoryName)
		if err != nil {
			return nil, err
		}
		categories = append(categories, *category)
	}

	noteByTitle, _ := ns.noteRepo.GetByTitle(ctx, formattedTitle)
	if noteByTitle != nil {
		return nil, validations.ErrDuplicateTitle
	}

	note := models.NewNote(formattedTitle, formattedContent, categories, userID)

	if _, err := ns.noteRepo.Create(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (ns *NoteService) GetNoteById(ctx context.Context, noteId uint) (*models.Note, error) {
	return ns.noteRepo.GetNoteById(ctx, noteId)
}

func (ns *NoteService) GetAllNotes(ctx context.Context) ([]*models.Note, error) {
	notes, err := ns.noteRepo.GetAllNotes(ctx)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (ns *NoteService) GetNotesByCategories(ctx context.Context, categoryNames []string) ([]*models.Note, error) {
	if len(categoryNames) == 0 {
		return nil, validations.ErrEmptyCategoryFilter
	}

	var formattedCategoryNames []string
	for _, name := range categoryNames {
		valid, formattedName, err := utils.ValidateAndFormatCategory(name)
		if !valid {
			return nil, err
		}
		formattedCategoryNames = append(formattedCategoryNames, formattedName)
	}

	notes, err := ns.noteRepo.GetNotesByCategories(ctx, formattedCategoryNames)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (ns *NoteService) UpdateNote(ctx context.Context, noteId uint, updatedNote *models.Note) (*uint, error) {

	existingNote, err := ns.GetNoteById(ctx, noteId)
	if err != nil {
		return nil, err
	}

	if existingNote.Categories == nil {
		existingNote.Categories = []models.Category{}
	}

	validTitle, formattedTitle, err := utils.ValidateAndFormatTitle(updatedNote.Title)
	if !validTitle {
		return nil, err
	}
	potentialDif, _ := ns.noteRepo.GetByTitle(ctx, formattedTitle)
	if potentialDif != nil && potentialDif.ID != existingNote.ID {
		return nil, validations.ErrDuplicateTitle
	}
	validContent, formattedContent, err := utils.ValidateAndFormatContent(updatedNote.Content)
	if !validContent {
		return nil, err
	}
	uniqueCats := make(map[string]models.Category)
	for _, c := range updatedNote.Categories {
		valid, name, err := utils.ValidateAndFormatCategory(c.Name)
		if !valid {
			return nil, err
		}
		nCat, err := ns.CategoryService.GetByNameOrCreate(ctx, name)
		if err != nil {
			return nil, err
		}
		uniqueCats[nCat.Name] = *nCat
	}

	newCats := make([]models.Category, 0, len(uniqueCats))
	for _, cat := range uniqueCats {
		newCats = append(newCats, cat)
	}

	if len(newCats) > 4 {
		return nil, validations.ErrTooManyCat
	}
	if len(newCats) < 1 {
		return nil, validations.ErrZeroCategory
	}

	if existingNote.Title == formattedTitle &&
		existingNote.Content == formattedContent &&
		existingNote.IsArchived == updatedNote.IsArchived &&
		utils.CompareCategories(existingNote.Categories, newCats) {
		return nil, validations.ErrNoChangesDetected
	}
	existingNote.Categories = newCats
	existingNote.Title = formattedTitle
	existingNote.Content = formattedContent
	existingNote.IsArchived = updatedNote.IsArchived
	existingNote.UpdatedAt = date.ArgentinaTimeNow()

	updatedId, err := ns.noteRepo.UpdateNote(ctx, existingNote)
	if err != nil {
		return nil, err
	}

	return updatedId, nil
}

func (ns *NoteService) DeleteNote(ctx context.Context, noteId uint) (*uint, error) {
	note, err := ns.GetNoteById(ctx, noteId)
	if err != nil {
		return nil, err
	}
	return ns.noteRepo.Delete(ctx, note)
}

func (ns *NoteService) AddCategoryToNote(ctx context.Context, noteId uint, categoryName string) (*models.Note, error) {
	note, err := ns.GetNoteById(ctx, noteId)
	if err != nil {
		return nil, err
	}
	if len(note.Categories) == 4 {
		return nil, validations.ErrFullCatCount
	}

	valid, formmatedCatname, err := utils.ValidateAndFormatCategory(categoryName)
	if !valid {
		return nil, err
	}
	for _, cat := range note.Categories {
		if cat.Name == formmatedCatname {
			return nil, validations.ErrCatAlreadyAdded
		}
	}
	category, err := ns.CategoryService.GetByNameOrCreate(ctx, formmatedCatname)
	if err != nil {
		return nil, err
	}
	note.Categories = append(note.Categories, *category)
	note.UpdatedAt = date.ArgentinaTimeNow()
	if _, err = ns.UpdateNote(ctx, noteId, note); err != nil {
		return nil, err
	}
	return note, nil
}

func (ns *NoteService) RemoveCategoryFromNote(ctx context.Context, noteId uint, categoryName string) (*models.Note, error) {

	note, err := ns.GetNoteById(ctx, noteId)
	if err != nil {
		return nil, err
	}
	if len(note.Categories) == 1 {
		return nil, validations.ErrMinCategory
	}

	valid, formmatedCatname, err := utils.ValidateAndFormatCategory(categoryName)
	if !valid {
		return nil, err
	}

	found := false
	var updatedCategories []models.Category

	for _, cat := range note.Categories {
		if cat.Name == formmatedCatname {
			found = true
			continue
		}
		updatedCategories = append(updatedCategories, cat)
	}

	if !found {
		return nil, validations.ErrCategoryNotFound
	}

	note.Categories = updatedCategories
	note.UpdatedAt = date.ArgentinaTimeNow()

	if _, err = ns.UpdateNote(ctx, noteId, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (ns *NoteService) ToggleArchiveStatus(ctx context.Context, noteId uint) (*models.Note, error) {
	note, err := ns.GetNoteById(ctx, noteId)
	if err != nil {
		return nil, err
	}
	note.IsArchived = !note.IsArchived
	note.UpdatedAt = date.ArgentinaTimeNow()
	if _, err := ns.UpdateNote(ctx, noteId, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (ns *NoteService) FilterNotes(ctx context.Context, isArchived *bool, categories []string) ([]models.Note, error) {

	var filteredCategories []string
	if len(categories) > 0 {
		filteredCategories = categories
	}
	notes, err := ns.noteRepo.FilterNotes(ctx, isArchived, filteredCategories)
	if err != nil {
		return nil, err
	}
	if len(notes) == 0 {
		return nil, validations.ErrNoNotesFound
	}
	return notes, nil
}

func (ns *NoteService) DeleteNoteForUser(ctx context.Context, noteId uint, userId uint) (*uint, error) {
	deletedNoteId, err := ns.noteRepo.DeleteNoteByUserId(ctx, noteId, userId)
	if err != nil {
		return nil, err
	}
	return deletedNoteId, nil
}
