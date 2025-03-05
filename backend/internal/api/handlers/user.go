package handlers

import (
	"net/http"
	"notes/internal/models"
	"notes/internal/services"
	"notes/pkg/request"
	"notes/pkg/response"
	"notes/pkg/utils"
	"notes/pkg/validations"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserService *services.UserService
	HttpErrs    *HttpErrors
}

func NewUserHandler(userService *services.UserService, httpErr *HttpErrors) *UserHandler {
	return &UserHandler{
		UserService: userService,
		HttpErrs:    httpErr,
	}
}

// RegisterUserHandler godoc
// @Summary     Register new user
// @Description Create user account with username and password
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       request body UserRequest true " "
// @Success     201 {object} UserResponse
// @Failure     400 {object} ErrorResponse
// @Failure     409 {object} ErrorResponse
// @Failure     401 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /user/register [post]
func (uh *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	err := request.DecodeJSONStrict(w, r, &req)
	if err != nil {
		uh.HttpErrs.badRequest(w, r, err, ReqErrKey)
		return
	}
	if len(req.Password) < 5 || len(req.Password) > 20 {
		uh.HttpErrs.CheckErrType(w, r, validations.ErrUserNamePassLength)
		return
	}

	user, err := uh.UserService.RegisterUser(r.Context(), w, req.UserName, req.Password)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}
	userResponse := UserResponse{
		ID:        &user.ID,
		UserName:  user.UserName,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
	}

	if err := response.JSON(w, http.StatusCreated, userResponse); err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}
}

// LoginUserHandler godoc
// @Summary     Authenticate user
// @Description Log in with username and password
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       credentials body UserRequest true " "
// @Success     200 {object} UserResponse
// @Failure     400 {object} ErrorResponse
// @Failure     401 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /user/login [post]
func (uh *UserHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	err := request.DecodeJSONStrict(w, r, &req)
	if err != nil {
		uh.HttpErrs.badRequest(w, r, err, ReqErrKey)
		return
	}
	if len(req.Password) < 5 || len(req.Password) > 20 {
		uh.HttpErrs.CheckErrType(w, r, validations.ErrUserNamePassLength)
		return
	}
	usr, err := uh.UserService.LoginUser(r.Context(), w, req.UserName, req.Password)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}
	if usr == nil {
		uh.HttpErrs.CheckErrType(w, r, validations.ErrInvalidCredentials)
		return
	}
	var res UserResponse
	res.ID = &usr.ID
	res.UserName = usr.UserName
	res.Notes = usr.Notes
	res.CreatedAt = usr.CreatedAt
	res.UpdatedAt = usr.UpdatedAt
	if err := response.JSON(w, http.StatusOK, &res); err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}
}

func (uh *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	user, err := uh.UserService.GetUserByID(r.Context(), *userID)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	userResponse := UserResponse{
		UserName:  user.UserName,
		Notes:     user.Notes,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	if err := response.JSON(w, http.StatusOK, userResponse); err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}
}

func (uh *UserHandler) AuthCheckHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, validations.ErrUnauthorized)
		return
	}

	res := map[string]interface{}{
		"authenticated": true,
		"user_id":       *userID,
	}
	if err := response.JSON(w, http.StatusOK, res); err != nil {
		uh.HttpErrs.CheckErrType(w, r, validations.ErrUnauthorized)
		return
	}
}

// Logout godoc
// @Summary     logout connected user
// @Description logout
// @Tags        users
// @Accept      json
// @Produce     json
// @Router      /user/logout [post]
func (uh *UserHandler) LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	err := uh.UserService.LogoutUser(w)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}
	if err := response.JSON(w, http.StatusOK, map[string]string{"message": "Logout successful"}); err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}
}

// CreateNoteHandler creates a new note for the user.
// @Summary Create a new note
// @Description Adds a new note for the authenticated user.
// @Tags notes
// @Security notes_jwt
// @Accept json
// @Produce json
// @Param note body CreateNoteRequest true "Note data"
// @Success 201 {object} APIResponse "Note created successfully"
// @Failure 400 {object} ErrorResponse "Invalid input data"
// @Failure 401 {object} ErrorResponse "Unauthorized access"
// @Router /notes [post]
func (uh *UserHandler) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {

	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}
	var req CreateNoteRequest
	err = request.DecodeJSONStrict(w, r, &req)
	if err != nil {
		uh.HttpErrs.badRequest(w, r, err, ReqErrKey)
		return
	}
	note, err := uh.UserService.CreateNote(r.Context(), req.Title, req.Content, req.Categories, *userID)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	noteResponse := GetNoteResponse{
		ID:         note.ID,
		Title:      note.Title,
		Content:    note.Content,
		Categories: note.Categories,
		IsArchived: note.IsArchived,
		CreatedAt:  note.CreatedAt,
		UpdatedAt:  note.UpdatedAt,
		UserID:     note.UserID,
	}

	if err := response.JSON(w, http.StatusCreated, noteResponse); err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}
}

// GetNoteByIdHandler retrieves a specific note by its ID for the authenticated user.
// @Summary Retrieve a specific note by ID
// @Description Fetches a single note by its unique ID for the authenticated user.
// @Tags notes
// @Security notes_jwt
// @Produce json
// @Param noteId path int true "Note ID"
// @Success 200 {object} GetNoteResponse "Note retrieved successfully"
// @Failure 400 {object} ErrorResponse "Invalid note ID"
// @Failure 401 {object} ErrorResponse "Unauthorized access"
// @Failure 404 {object} ErrorResponse "Note not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /notes/{noteId} [get]
func (uh *UserHandler) GetNoteByIdHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	noteID, err := strconv.ParseUint(vars["noteId"], 10, 64)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, validations.ErrInlvalidId)
		return
	}

	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	note, err := uh.UserService.GetNoteById(r.Context(), uint(noteID), *userID)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	noteResponse := GetNoteResponse{
		ID:         note.ID,
		Title:      note.Title,
		Content:    note.Content,
		Categories: note.Categories,
		IsArchived: note.IsArchived,
		CreatedAt:  note.CreatedAt,
		UpdatedAt:  note.UpdatedAt,
		UserID:     note.UserID,
	}

	if err := response.JSON(w, http.StatusOK, noteResponse); err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
	}
}

// UpdateNoteHandler updates a specific note for the authenticated user.
// @Summary Update a specific note by ID
// @Description Updates the note for the authenticated user by its unique ID.
// @Tags notes
// @Security notes_jwt
// @Accept json
// @Produce json
// @Param noteId path int true "Note ID"
// @Param note body UpdateNoteRequest true "Updated note data"
// @Success 200 {object} APIResponse "Note ID of updated note"
// @Failure 400 {object} ErrorResponse "Invalid note ID or request data"
// @Failure 401 {object} ErrorResponse "Unauthorized access"
// @Failure 404 {object} ErrorResponse "Note not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /notes/{noteId} [put]
func (uh *UserHandler) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteID, err := strconv.ParseUint(vars["noteId"], 10, 32)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, validations.ErrInlvalidId)
		return
	}

	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	var req UpdateNoteRequest
	err = request.DecodeJSONStrict(w, r, &req)
	if err != nil {
		uh.HttpErrs.badRequest(w, r, err, ReqErrKey)
		return
	}

	updated := models.NewNote(req.Title, req.Content, req.Categories, *userID)
	updatedNoteId, err := uh.UserService.UpdateNoteForUser(r.Context(), *userID, uint(noteID), updated)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}
	if err := response.JSON(w, http.StatusOK, updatedNoteId); err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
	}
}

// DeleteNoteHandler deletes a specific note by its ID for the authenticated user.
// @Summary Delete a specific note by ID
// @Description Deletes a note by its unique ID for the authenticated user.
// @Tags notes
// @Security notes_jwt
// @Accept json
// @Produce json
// @Param noteId path int true "Note ID"
// @Success 200 {object} APIResponse "Note ID of the deleted note"
// @Failure 400 {object} ErrorResponse "Invalid note ID"
// @Failure 401 {object} ErrorResponse "Unauthorized access"
// @Failure 404 {object} ErrorResponse "Note not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /notes/{noteId} [delete]
func (uh *UserHandler) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteID, err := strconv.ParseUint(vars["noteId"], 10, 32)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, validations.ErrInlvalidId)
		return
	}
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	_, err = uh.UserService.DeleteNoteForUser(r.Context(), *userID, uint(noteID))
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, noteID); err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
	}
}

// GetAllNotesByUserHandler retrieves all notes for the authenticated user.
// @Summary Retrieve all notes
// @Description Fetches all notes created by the authenticated user.
// @Tags notes
// @Security notes_jwt
// @Produce json
// @Success 200 {array} GetNoteResponse "List of notes"
// @Failure 401 {object} ErrorResponse "Unauthorized access"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /notes [get]
func (uh *UserHandler) GetAllNotesByUserHandler(w http.ResponseWriter, r *http.Request) {

	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	notes, err := uh.UserService.GetAllNotesByUserID(r.Context(), *userID)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	noteResponses := make([]GetNoteResponse, 0, len(notes))
	for _, note := range notes {
		noteResponses = append(noteResponses, GetNoteResponse{
			ID:         note.ID,
			Title:      note.Title,
			Content:    note.Content,
			Categories: note.Categories,
			IsArchived: note.IsArchived,
			CreatedAt:  note.CreatedAt,
			UpdatedAt:  note.UpdatedAt,
			UserID:     note.UserID,
		})
	}

	if err := response.JSON(w, http.StatusOK, noteResponses); err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
	}
}

// AddCategoryToNoteHandler adds a category to a specific note for the authenticated user.
// @Summary Add a category to a specific note by ID
// @Description Adds a category to a note for the authenticated user.
// @Tags notes
// @Security notes_jwt
// @Accept json
// @Produce json
// @Param noteId path int true "Note ID"
// @Param categoryName path string true "Category name"
// @Success 200 {object} GetNoteResponse "Updated note with added category"
// @Failure 400 {object} ErrorResponse "Invalid note ID or category name"
// @Failure 401 {object} ErrorResponse "Unauthorized access"
// @Failure 404 {object} ErrorResponse "Note or category not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /notes/{noteId}/categories/{categoryName} [post]
func (uh *UserHandler) AddCategoryToNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteID, err := strconv.ParseUint(vars["noteId"], 10, 32)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, validations.ErrInlvalidId)
		return
	}

	categoryName := vars["categoryName"]
	if categoryName == "" {
		uh.HttpErrs.CheckErrType(w, r, validations.ErrCategoryName)
		return
	}
	valid, formmatedCat, err := utils.ValidateAndFormatCategory(categoryName)
	if !valid {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	updatedNote, err := uh.UserService.AddCategoryToNoteForUser(r.Context(), *userID, uint(noteID), formmatedCat)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	resp := GetNoteResponse{
		ID:         updatedNote.ID,
		Title:      updatedNote.Title,
		Content:    updatedNote.Content,
		Categories: updatedNote.Categories,
		IsArchived: updatedNote.IsArchived,
		CreatedAt:  updatedNote.CreatedAt,
		UpdatedAt:  updatedNote.UpdatedAt,
		UserID:     updatedNote.UserID,
	}

	if err := response.JSON(w, http.StatusOK, resp); err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
	}
}

// RemoveCategoryFromNoteHandler removes a category from a specific note for the authenticated user.
// @Summary Remove a category from a specific note by ID
// @Description Removes a category from a note for the authenticated user.
// @Tags notes
// @Security notes_jwt
// @Accept json
// @Produce json
// @Param noteId path int true "Note ID"
// @Param categoryName path string true "Category name"
// @Success 200 {object} GetNoteResponse "Updated note with removed category"
// @Failure 400 {object} ErrorResponse "Invalid note ID or category name"
// @Failure 401 {object} ErrorResponse "Unauthorized access"
// @Failure 404 {object} ErrorResponse "Note or category not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /notes/{noteId}/categories/{categoryName} [delete]
func (uh *UserHandler) RemoveCategoryFromNoteHandler(w http.ResponseWriter, r *http.Request) {

	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	vars := mux.Vars(r)
	noteID, err := strconv.ParseUint(vars["noteId"], 10, 32)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, validations.ErrInlvalidId)
		return
	}

	categoryName := vars["categoryName"]
	if categoryName == "" {
		uh.HttpErrs.CheckErrType(w, r, validations.ErrCategoryName)
		return
	}
	valid, formmatedCat, err := utils.ValidateAndFormatCategory(categoryName)
	if !valid {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	updatedNote, err := uh.UserService.RemoveCategoryFromNoteForUser(r.Context(), *userID, uint(noteID), formmatedCat)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	resp := GetNoteResponse{
		ID:         updatedNote.ID,
		Title:      updatedNote.Title,
		Content:    updatedNote.Content,
		Categories: updatedNote.Categories,
		IsArchived: updatedNote.IsArchived,
		CreatedAt:  updatedNote.CreatedAt,
		UpdatedAt:  updatedNote.UpdatedAt,
		UserID:     updatedNote.UserID,
	}

	if err := response.JSON(w, http.StatusOK, resp); err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
	}
}

// ToggleArchiveStatusHandler toggles the archive status of a specific note for the authenticated user.
// @Summary Toggle the archive status of a note
// @Description Toggles whether a note is archived or not for the authenticated user.
// @Tags notes
// @Security notes_jwt
// @Accept json
// @Produce json
// @Param noteId path int true "Note ID"
// @Success 200 {object} GetNoteResponse "Note with updated archive status"
// @Failure 400 {object} ErrorResponse "Invalid note ID"
// @Failure 401 {object} ErrorResponse "Unauthorized access"
// @Failure 404 {object} ErrorResponse "Note not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /notes/{noteId}/archive [put]
func (uh *UserHandler) ToggleArchiveStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteID, err := strconv.ParseUint(vars["noteId"], 10, 64)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, validations.ErrInlvalidId)
		return
	}

	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	updatedNote, err := uh.UserService.ToggleArchiveStatusForUser(r.Context(), *userID, uint(noteID))
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	resp := GetNoteResponse{
		ID:         updatedNote.ID,
		Title:      updatedNote.Title,
		Content:    updatedNote.Content,
		Categories: updatedNote.Categories,
		IsArchived: updatedNote.IsArchived,
		CreatedAt:  updatedNote.CreatedAt,
		UpdatedAt:  updatedNote.UpdatedAt,
		UserID:     updatedNote.UserID,
	}

	if err := response.JSON(w, http.StatusOK, resp); err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
	}
}

// FilterNotesForUserHandler filters notes by categories and archived status.
// @Summary Filter notes by categories and archived status
// @Description Filters notes based on categories and archived status. Both parameters are optional.
// @Tags notes
// @Security notes_jwt
// @Param isArchived query bool false "Filter by archived status (optional)"
// @Param categories query []string false "Filter notes by categories (optional)" "List of categories"
// @Success 200 {array} GetNoteResponse "Filtered notes"
// @Failure 400 {object} ErrorResponse "Invalid query parameters"
// @Failure 401 {object} ErrorResponse "Unauthorized access"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /notes/filter [get]
func (uh *UserHandler) FilterNotesForUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	isArchived := r.URL.Query().Get("isArchived")
	var archived *bool
	if isArchived != "" {
		value, err := strconv.ParseBool(isArchived)
		if err == nil {
			archived = &value
		}
	}

	categories := r.URL.Query()["categories"]

	notes, err := uh.UserService.FilterNotesForUser(r.Context(), *userID, archived, categories)
	if err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	noteResponses := make([]GetNoteResponse, 0, len(notes))
	for _, note := range notes {
		noteResponses = append(noteResponses, GetNoteResponse{
			ID:         note.ID,
			Title:      note.Title,
			Content:    note.Content,
			Categories: note.Categories,
			IsArchived: note.IsArchived,
			CreatedAt:  note.CreatedAt,
			UpdatedAt:  note.UpdatedAt,
			UserID:     note.UserID,
		})
	}

	if err := response.JSON(w, http.StatusOK, noteResponses); err != nil {
		uh.HttpErrs.CheckErrType(w, r, err)
	}
}
