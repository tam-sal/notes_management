package handlers

import "notes/internal/services"

type CategoryHandler struct {
	CatService *services.CategoryService
	HttpErrs   *HttpErrors
}

func NewCategoryHandler(cs *services.CategoryService, httpErrs *HttpErrors) *CategoryHandler {
	return &CategoryHandler{
		CatService: cs,
		HttpErrs:   httpErrs,
	}
}
