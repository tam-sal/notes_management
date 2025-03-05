package handlers

import (
	"net/http"
	"notes/internal/models"
	"notes/internal/services"
	"notes/pkg/request"
	"notes/pkg/response"
	"notes/pkg/validations"
	"strconv"

	"github.com/gorilla/mux"
)

type NoteHandler struct {
	NoteService *services.NoteService
	HttpErrs    *HttpErrors
}

func NewNoteHandler(noteService *services.NoteService, httpErr *HttpErrors) *NoteHandler {
	return &NoteHandler{
		NoteService: noteService,
		HttpErrs:    httpErr,
	}
}

func (nh *NoteHandler) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	var req *CreateNoteRequest
	err := request.DecodeJSONStrict(w, r, &req)
	if err != nil {
		nh.HttpErrs.badRequest(w, r, err, ReqErrKey)
	}

	note, err := nh.NoteService.CreateNote(r.Context(), req.Title, req.Content, req.Categories, req.UserID)
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}
	if err := response.JSON(w, http.StatusOK, note); err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}
}

func (nh *NoteHandler) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if id == "" || !ok {
		nh.HttpErrs.CheckErrType(w, r, validations.ErrMissingId)
		return
	}

	intId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, validations.ErrInlvalidId)
		return
	}
	uintID := uint(intId)

	var req UpdateNoteRequest
	err = request.DecodeJSONStrict(w, r, &req)
	if err != nil {
		nh.HttpErrs.badRequest(w, r, err, ReqErrKey)
		return
	}

	if len(req.Categories) == 0 {
		nh.HttpErrs.badRequest(w, r, validations.ErrEmptyCategory, ReqErrKey)
		return
	}

	updatedNote := models.Note{
		ID:         uintID,
		Title:      req.Title,
		Content:    req.Content,
		Categories: req.Categories,
		IsArchived: req.IsArchived,
		UserID:     req.UserID,
	}

	ider, err := nh.NoteService.UpdateNote(r.Context(), uintID, &updatedNote)
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	if ider == nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	if err := response.JSON(w, http.StatusOK, ider); err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}
}

func (nh *NoteHandler) GetNoteByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if id == "" || !ok {
		nh.HttpErrs.CheckErrType(w, r, validations.ErrMissingId)
		return
	}
	intId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, validations.ErrInlvalidId)
		return
	}
	uintID := uint(intId)

	note, err := nh.NoteService.GetNoteById(r.Context(), uintID)
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	res := &GetNoteResponse{
		ID:         note.ID,
		Title:      note.Title,
		Content:    note.Content,
		Categories: note.Categories,
		UserID:     note.UserID,
		CreatedAt:  note.CreatedAt,
		UpdatedAt:  note.UpdatedAt,
	}

	if err = response.JSON(w, http.StatusOK, res); err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}

}

func (nh *NoteHandler) GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := nh.NoteService.GetAllNotes(r.Context())
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}
	var res []*GetNoteResponse
	for _, n := range notes {
		rs := GetNoteResponse{
			ID:         n.ID,
			Title:      n.Title,
			Content:    n.Content,
			IsArchived: n.IsArchived,
			Categories: n.Categories,
			UserID:     n.UserID,
			CreatedAt:  n.CreatedAt,
			UpdatedAt:  n.UpdatedAt,
		}
		res = append(res, &rs)
	}
	if err := response.JSON(w, http.StatusOK, res); err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}
}

func (nh *NoteHandler) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {

	id, ok := mux.Vars(r)["id"]
	if id == "" || !ok {
		nh.HttpErrs.CheckErrType(w, r, validations.ErrMissingId)
		return
	}
	intId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, validations.ErrInlvalidId)
		return
	}
	uintID := uint(intId)
	uid, err := nh.NoteService.DeleteNote(r.Context(), uintID)
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}
	if err := response.JSON(w, http.StatusOK, uid); err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}
}

func (nh *NoteHandler) RemoveCategoryFromNoteHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	noteIdStr, noteIdExists := vars["noteId"]
	categoryName, catNameExists := vars["categoryName"]

	if noteIdStr == "" || !noteIdExists || !catNameExists || categoryName == "" {
		nh.HttpErrs.CheckErrType(w, r, validations.ErrMissingParameters)
		return
	}

	noteId, err := strconv.ParseUint(noteIdStr, 10, 32)
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, validations.ErrInlvalidId)
		return
	}

	updatedNote, err := nh.NoteService.RemoveCategoryFromNote(r.Context(), uint(noteId), categoryName)
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}
	res := &GetNoteResponse{
		ID:         updatedNote.ID,
		Title:      updatedNote.Title,
		Content:    updatedNote.Content,
		Categories: updatedNote.Categories,
		UserID:     updatedNote.UserID,
		CreatedAt:  updatedNote.CreatedAt,
		UpdatedAt:  updatedNote.UpdatedAt,
	}

	if err := response.JSON(w, http.StatusOK, &res); err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}
}

func (nh *NoteHandler) AddCategoryToNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, idOk := vars["noteId"]
	category, categoryOk := vars["categoryName"]

	if !idOk || !categoryOk || id == "" || category == "" {
		nh.HttpErrs.CheckErrType(w, r, validations.ErrMissingParameters)
		return
	}

	intId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, validations.ErrInlvalidId)
		return
	}
	uintID := uint(intId)

	note, err := nh.NoteService.AddCategoryToNote(r.Context(), uintID, category)
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}

	res := &GetNoteResponse{
		ID:         note.ID,
		Title:      note.Title,
		Content:    note.Content,
		Categories: note.Categories,
		UserID:     note.UserID,
		CreatedAt:  note.CreatedAt,
		UpdatedAt:  note.UpdatedAt,
	}

	if err := response.JSON(w, http.StatusOK, &res); err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}
}

func (nh *NoteHandler) ToggleArchiveStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, idOk := vars["noteId"]

	if !idOk || id == "" {
		nh.HttpErrs.CheckErrType(w, r, validations.ErrMissingId)
		return
	}

	intId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, validations.ErrInlvalidId)
		return
	}
	uintID := uint(intId)

	note, err := nh.NoteService.ToggleArchiveStatus(r.Context(), uintID)
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}
	res := &GetNoteResponse{
		ID:         note.ID,
		Title:      note.Title,
		Content:    note.Content,
		Categories: note.Categories,
		UserID:     note.UserID,
		CreatedAt:  note.CreatedAt,
		UpdatedAt:  note.UpdatedAt,
	}

	if err := response.JSON(w, http.StatusOK, &res); err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}
}

func (nh *NoteHandler) FilterNotesHandler(w http.ResponseWriter, r *http.Request) {

	isArchivedParam := r.URL.Query().Get("isArchived")
	categoriesParam := r.URL.Query()["categories"]

	var isArchived *bool
	if isArchivedParam != "" {
		archivedBool, err := strconv.ParseBool(isArchivedParam)
		if err != nil {
			nh.HttpErrs.CheckErrType(w, r, validations.ErrInvalidArchivedValue)
			return
		}
		isArchived = &archivedBool
	}

	notes, err := nh.NoteService.FilterNotes(r.Context(), isArchived, categoriesParam)
	if err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}
	var res []*GetNoteResponse
	for _, n := range notes {
		rs := GetNoteResponse{
			ID:         n.ID,
			Title:      n.Title,
			Content:    n.Content,
			IsArchived: n.IsArchived,
			Categories: n.Categories,
			UserID:     n.UserID,
			CreatedAt:  n.CreatedAt,
			UpdatedAt:  n.UpdatedAt,
		}
		res = append(res, &rs)
	}
	if err := response.JSON(w, http.StatusOK, &res); err != nil {
		nh.HttpErrs.CheckErrType(w, r, err)
		return
	}
}
