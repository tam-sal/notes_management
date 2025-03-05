package main

import (
	"net/http"
	_ "notes/cmd/api/docs"
	"notes/internal/api/handlers"
	"time"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	mx := mux.NewRouter()
	app.logger.Debug("INTIALIZING ROUTES")
	const GET, POST, PUT, DELETE, OPTIONS = "GET", "POST", "PUT", "DELETE", "OPTIONS"

	// Rate Limiter
	rateLimiter := handlers.NewRateLimiter(5, 10)

	// Middlewares
	mx.Use(app.handlers.RecoverPanic)
	mx.Use(app.handlers.AddHeadersWithCSP)
	mx.Use(app.handlers.MetricsMiddleware)
	mx.Use(app.handlers.WithTimeout(time.Second * 20))
	mx.Use(app.handlers.LimitMiddleware(rateLimiter))
	mx.Use(app.handlers.LogAccess)

	// Custom Handling
	mx.NotFoundHandler = http.HandlerFunc(app.handlers.HttpErrs.NotFound)
	mx.MethodNotAllowedHandler = http.HandlerFunc(app.handlers.HttpErrs.MethodNotAllowed)

	// Swagger Documentation Route
	// Register Swagger documentation route
	// @securityDefinitions.apikey notes_jwt
	// @in cookies
	// @name notes_jwt
	app.handlers.RegisterSwaggerHandler(mx)

	// Routes
	// Metrics
	mx.HandleFunc("/prometheus-metrics", app.handlers.MetricsHandler).Methods("GET").Name("metrics")
	// Health Check
	mx.HandleFunc("/status", app.handlers.StatusHandler).Methods(GET, OPTIONS).Name("status")

	// SUBROUTES

	userRouter := mx.PathPrefix("/user").Subrouter()
	noteRouter := mx.PathPrefix("/notes").Subrouter()

	// USER ROUTES
	userRouter.HandleFunc("/register", app.handlers.UserHandler.RegisterUserHandler).Methods(POST, OPTIONS).Name("user:register")
	userRouter.HandleFunc("/login", app.handlers.UserHandler.LoginUserHandler).Methods(POST, OPTIONS).Name("user:login")
	userRouter.HandleFunc("/logout", app.handlers.UserHandler.LogoutUserHandler).Methods(POST, OPTIONS).Name("user:logout")
	userRouter.Handle("/auth-check", app.handlers.PROTECT(http.HandlerFunc(app.handlers.UserHandler.AuthCheckHandler))).Methods(GET, OPTIONS).Name("user:auth-check")

	// NOTES (PROTECTED) ROUTES
	noteRouter.Use(app.handlers.PROTECT)
	noteRouter.HandleFunc("/filter", app.handlers.UserHandler.FilterNotesForUserHandler).Methods(GET, OPTIONS).Name("notes:filter")
	noteRouter.HandleFunc("", app.handlers.UserHandler.GetAllNotesByUserHandler).Methods(GET, OPTIONS).Name("notes:list")
	noteRouter.HandleFunc("", app.handlers.UserHandler.CreateNoteHandler).Methods(POST, OPTIONS).Name("notes:create")
	noteRouter.HandleFunc("/{noteId}", app.handlers.UserHandler.GetNoteByIdHandler).Methods(GET, OPTIONS).Name("notes:get")
	noteRouter.HandleFunc("/{noteId}", app.handlers.UserHandler.UpdateNoteHandler).Methods(PUT, OPTIONS).Name("notes:update")
	noteRouter.HandleFunc("/{noteId}", app.handlers.UserHandler.DeleteNoteHandler).Methods(DELETE, OPTIONS).Name("notes:delete")
	noteRouter.HandleFunc("/{noteId}/archive-toggle", app.handlers.UserHandler.ToggleArchiveStatusHandler).Methods(PUT, OPTIONS).Name("archive-toggle")

	noteRouter.HandleFunc("/{noteId}/categories/{categoryName}", app.handlers.UserHandler.AddCategoryToNoteHandler).Methods(POST, OPTIONS).Name("category:add")
	noteRouter.HandleFunc("/{noteId}/categories/{categoryName}", app.handlers.UserHandler.RemoveCategoryFromNoteHandler).Methods(DELETE, OPTIONS).Name("category:remove")

	// Default Fallback Handling
	mx.PathPrefix("/").HandlerFunc(app.handlers.HttpErrs.NotFound).Name("fallback")

	return mx
}
