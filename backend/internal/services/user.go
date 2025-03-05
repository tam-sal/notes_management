package services

import (
	"context"
	"net/http"

	"notes/internal/models"
	"notes/internal/repositories"
	"notes/pkg/utils"
	"notes/pkg/validations"
	"time"
)

type UserService struct {
	userRepo    *repositories.UserRepository
	noteService *NoteService
}

func NewUserService(userRepo *repositories.UserRepository, noteService *NoteService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		noteService: noteService,
	}
}

func (us *UserService) CreateUser(ctx context.Context, username, password string) (*models.User, error) {

	valid, formattedUsername, err := utils.ValidateAndFormatUsername(username)
	if !valid {
		return nil, err
	}

	_, err = us.GetUserByUsername(ctx, *formattedUsername)
	if err == nil {
		return nil, validations.ErrUserAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.NewUser(*formattedUsername, *hashedPassword, nil)

	if err := us.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	return us.userRepo.GetUserByID(ctx, userID)
}

func (us *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return us.userRepo.GetUserByUsername(ctx, username)
}

func (us *UserService) AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	user, err := us.GetUserByUsername(ctx, username)
	if err != nil || user == nil {
		return nil, validations.ErrInvalidCredentials
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, validations.ErrInvalidCredentials
	}
	return user, nil
}

func (us *UserService) CreateNote(ctx context.Context, title, content string, categoryNames []string, userID uint) (*models.Note, error) {
	note, err := us.noteService.CreateNote(ctx, title, content, categoryNames, userID)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (us *UserService) GetNoteById(ctx context.Context, noteId uint, userId uint) (*models.Note, error) {

	note, err := us.noteService.GetNoteById(ctx, noteId)
	if err != nil {
		return nil, err
	}
	if note.UserID != userId {
		return nil, validations.ErrNoteNotOwnedByUser
	}
	return note, nil
}

func (us *UserService) UpdateNoteForUser(ctx context.Context, userId uint, noteId uint, updatedNote *models.Note) (*uint, error) {
	existingNote, err := us.noteService.GetNoteById(ctx, noteId)
	if err != nil {
		return nil, err
	}

	if existingNote.UserID != userId {
		return nil, validations.ErrNoteNotOwnedByUser
	}

	updatedNoteId, err := us.noteService.UpdateNote(ctx, noteId, updatedNote)
	if err != nil {
		return nil, err
	}

	return updatedNoteId, nil
}

func (us *UserService) DeleteNoteForUser(ctx context.Context, userId uint, noteId uint) (*uint, error) {

	id, err := us.noteService.DeleteNoteForUser(ctx, noteId, userId)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (us *UserService) GetAllNotesByUserID(ctx context.Context, userId uint) ([]models.Note, error) {
	user, err := us.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	if len(user.Notes) > 0 {
		return user.Notes, nil
	}
	return nil, validations.ErrNoNotesFound
}

func (us *UserService) AddCategoryToNoteForUser(ctx context.Context, userId uint, noteId uint, categoryName string) (*models.Note, error) {
	note, err := us.noteService.GetNoteById(ctx, noteId)
	if err != nil {
		return nil, err
	}
	if note.UserID != userId {
		return nil, validations.ErrNoteNotOwnedByUser
	}
	updatedNote, err := us.noteService.AddCategoryToNote(ctx, noteId, categoryName)
	if err != nil {
		return nil, err
	}
	return updatedNote, nil
}

func (us *UserService) RemoveCategoryFromNoteForUser(ctx context.Context, userId uint, noteId uint, categoryName string) (*models.Note, error) {
	note, err := us.noteService.GetNoteById(ctx, noteId)
	if err != nil {
		return nil, err
	}
	if note.UserID != userId {
		return nil, validations.ErrNoteNotOwnedByUser
	}
	updatedNote, err := us.noteService.RemoveCategoryFromNote(ctx, noteId, categoryName)
	if err != nil {
		return nil, err
	}
	return updatedNote, nil
}

func (us *UserService) ToggleArchiveStatusForUser(ctx context.Context, userId uint, noteId uint) (*models.Note, error) {
	note, err := us.noteService.GetNoteById(ctx, noteId)
	if err != nil {
		return nil, err
	}
	if note.UserID != userId {
		return nil, validations.ErrNoteNotOwnedByUser
	}
	updatedNote, err := us.noteService.ToggleArchiveStatus(ctx, noteId)
	if err != nil {
		return nil, err
	}
	return updatedNote, nil
}

func (us *UserService) FilterNotesForUser(ctx context.Context, userId uint, isArchived *bool, categories []string) ([]models.Note, error) {

	notes, err := us.noteService.FilterNotes(ctx, isArchived, categories)
	if err != nil {
		return nil, err
	}

	var userNotes []models.Note
	for _, note := range notes {
		if note.UserID == userId {
			userNotes = append(userNotes, note)
		}
	}

	if len(userNotes) == 0 {
		return nil, validations.ErrNoNotesFound
	}

	return userNotes, nil
}

// Regular User
func (us *UserService) RegisterUser(ctx context.Context, w http.ResponseWriter, username, password string) (*models.User, error) {
	user, err := us.CreateUser(ctx, username, password)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return nil, validations.ErrTokenGeneration
	}
	utils.SetJWTAsCookie(w, token)
	return user, nil

}

func (us *UserService) LoginUser(ctx context.Context, w http.ResponseWriter, username, password string) (*models.User, error) {
	user, err := us.AuthenticateUser(ctx, username, password)
	if err != nil {
		return nil, err
	}
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return nil, validations.ErrTokenGeneration
	}
	utils.SetJWTAsCookie(w, token)
	return user, nil
}

// REFACTORED
func (us *UserService) LogoutUser(w http.ResponseWriter) error {

	http.SetCookie(w, &http.Cookie{
		Name:     "notes_jwt",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteNoneMode,
	})
	return nil
}
