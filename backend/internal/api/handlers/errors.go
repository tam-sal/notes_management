package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"notes/pkg/response"
	"notes/pkg/validations"
	"runtime"
	"runtime/debug"
	"strings"
)

type CustomErrKey string

const (
	APIErrKey CustomErrKey = "API"
	ReqErrKey CustomErrKey = "REQUEST"
	DBErrKey  CustomErrKey = "DATABASE"
	AUTH      CustomErrKey = "AUTH"
)

type HttpErrors struct {
	logger *slog.Logger
}

func NewHttpErrors(logger *slog.Logger) *HttpErrors {
	return &HttpErrors{
		logger: logger,
	}
}

func (h *HttpErrors) reportServerError(r *http.Request, err error) {
	message := err.Error()
	method := r.Method
	url := r.URL.String()
	trace := string(debug.Stack())

	shortTrace := strings.Join(strings.Split(trace, "\n")[:10], "\n")

	_, file, line, ok := runtime.Caller(5)
	if !ok {
		file = "unknown"
		line = 0
	}

	h.logger.Error(message,
		slog.Group("request",
			slog.String("method", method),
			slog.String("url", url),
		),
		slog.String("file", file),
		slog.Int("line", line),
		slog.String("trace", shortTrace),
	)
}

func (h *HttpErrors) errorMessage(w http.ResponseWriter, r *http.Request, status int, key CustomErrKey, message string, headers http.Header) {
	errObject := map[CustomErrKey]string{
		key: message,
	}
	err := response.ErrJSONWithHeaders(w, status, errObject, headers)
	if err != nil {
		h.reportServerError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *HttpErrors) ServerError(w http.ResponseWriter, r *http.Request, err error, key CustomErrKey) {
	h.reportServerError(r, err)
	h.errorMessage(w, r, http.StatusInternalServerError, key, err.Error(), nil)
}

// func (h *HttpErrors) notFoundResponse(w http.ResponseWriter, r *http.Request, key CustomErrKey, message string) {
// 	h.errorMessage(w, r, http.StatusNotFound, key, message, nil)
// }

// API
//
//	func (h *HttpErrors) notFound(w http.ResponseWriter, r *http.Request) {
//		msg := "the requested resource couldn't be found"
//		h.errorMessage(w, r, http.StatusNotFound, APIErrKey, msg, nil)
//	}
func (h *HttpErrors) NotFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	h.errorMessage(w, r, http.StatusNotFound, APIErrKey, message, nil)
}

func (h *HttpErrors) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	h.errorMessage(w, r, http.StatusMethodNotAllowed, APIErrKey, message, nil)
}

// REQUEST/client
// func (h *HttpErrors) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
// 	msg := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
// 	h.errorMessage(w, r, http.StatusBadRequest, ReqErrKey, msg, nil)
// }

// REQUEST/CLIENT - TOKEN
func (h *HttpErrors) Unauthorized(w http.ResponseWriter, r *http.Request, key CustomErrKey, message string) {
	h.errorMessage(w, r, http.StatusUnauthorized, key, message, nil)
}

// REQUEST/CLIENT - DB
func (h *HttpErrors) badRequest(w http.ResponseWriter, r *http.Request, err error, key CustomErrKey) {
	h.errorMessage(w, r, http.StatusBadRequest, key, err.Error(), nil)
}

// REQUEST/CLIENT - API
func (h *HttpErrors) gatewayTimeout(w http.ResponseWriter, r *http.Request, err error, key CustomErrKey) {
	h.errorMessage(w, r, http.StatusGatewayTimeout, key, err.Error(), nil)
}

func (h *HttpErrors) CheckErrType(w http.ResponseWriter, r *http.Request, err error) {

	if err == nil {
		return
	}

	switch {

	// REQUEST - CLIENT
	case errors.Is(err, validations.ErrCharactersContentExcess),
		errors.Is(err, validations.ErrCharactersExcessCat),
		errors.Is(err, validations.ErrEmptyContent),
		errors.Is(err, validations.ErrMissingId),
		errors.Is(err, validations.ErrInlvalidId),
		errors.Is(err, validations.ErrCharactersExcess),
		errors.Is(err, validations.ErreEmptyTitle),
		errors.Is(err, validations.ErrDuplicateTitle),
		errors.Is(err, validations.ErrCategoryNotFound),
		errors.Is(err, validations.ErrNoChangesDetected),
		errors.Is(err, validations.ErrFullCatCount),
		errors.Is(err, validations.ErrCatAlreadyAdded),
		errors.Is(err, validations.ErrEmptyCategory),
		errors.Is(err, validations.ErrMissingId),
		errors.Is(err, validations.ErrInlvalidId),
		errors.Is(err, validations.ErrMinCategory),
		errors.Is(err, validations.ErrMissingParameters),
		errors.Is(err, validations.ErrTitleFormat),
		errors.Is(err, validations.ErrNoteNotOwnedByUser),
		errors.Is(err, validations.ErrUserNamePassLength),
		errors.Is(err, validations.ErrUserAlreadyExists),
		errors.Is(err, validations.ErrTooManyCategories),
		errors.Is(err, validations.ErrRepeatedLetters),
		errors.Is(err, validations.ErrNoNotesFound),
		errors.Is(err, validations.ErreEmptyTitle),
		errors.Is(err, validations.ErrTooManyCat):
		h.badRequest(w, r, err, ReqErrKey)
		return

	// DATABASE
	case errors.Is(err, validations.ErrFetchingCategory),
		errors.Is(err, validations.ErrUserIdNotSet),
		errors.Is(err, validations.ErrFetchingNotes),
		errors.Is(err, validations.ErrFilterDB),
		errors.Is(err, validations.ErrCatAlreadyExist),
		errors.Is(err, validations.ErrCatCreate),
		errors.Is(err, validations.ErrCatNotFound),
		errors.Is(err, validations.ErrCatUpdate),
		errors.Is(err, validations.ErrCatDelete),
		errors.Is(err, validations.ErrNotTitle),
		errors.Is(err, validations.ErrFetchingCategories),
		errors.Is(err, validations.ErrNoteCreate),
		errors.Is(err, validations.ErrNoteNotFound),
		errors.Is(err, validations.ErrFetchingNote),
		errors.Is(err, validations.ErrAddNewCatToNote),
		errors.Is(err, validations.ErrNoteUpdate),
		errors.Is(err, validations.ErrNoteDelete):
		h.ServerError(w, r, err, DBErrKey)

	// API
	case errors.Is(err, context.Canceled),
		errors.Is(err, validations.ErrJsonResponse):
		h.ServerError(w, r, err, APIErrKey)

	case errors.Is(err, validations.ErrNotFound):
		h.ServerError(w, r, validations.ErrNotFound, APIErrKey)
	case errors.Is(err, validations.ErrRateLimitExcess):
		h.ServerError(w, r, validations.ErrRateLimitExcess, APIErrKey)

	// GATEWAY
	case errors.Is(err, context.DeadlineExceeded):
		h.gatewayTimeout(w, r, err, APIErrKey)

		//AUTH
	case errors.Is(err, validations.ErrHashingPwd):
		h.Unauthorized(w, r, AUTH, validations.ErrHashingPwd.Error())
	case errors.Is(err, validations.ErrInvalidCredentials):
		h.Unauthorized(w, r, AUTH, validations.ErrInvalidCredentials.Error())
	case errors.Is(err, validations.ErrTokenGeneration):
		h.Unauthorized(w, r, AUTH, validations.ErrTokenGeneration.Error())
	case errors.Is(err, validations.ErrJWT):
		h.Unauthorized(w, r, AUTH, validations.ErrJWT.Error())
	case errors.Is(err, validations.ErrTokenExpired):
		h.Unauthorized(w, r, AUTH, validations.ErrTokenExpired.Error())
	case errors.Is(err, validations.ErrTokenExpiry):
		h.Unauthorized(w, r, AUTH, validations.ErrTokenExpiry.Error())
	case errors.Is(err, validations.ErrInvalidUserID):
		h.Unauthorized(w, r, AUTH, validations.ErrInvalidUserID.Error())
	case errors.Is(err, validations.ErrUnauthorized):
		h.Unauthorized(w, r, AUTH, validations.ErrUnauthorized.Error())

	// DEFAULT
	default:
		h.ServerError(w, r, err, APIErrKey)
	}

}
