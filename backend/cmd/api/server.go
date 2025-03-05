package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"notes/internal/api/handlers"
	"notes/internal/configs"
	"notes/internal/db"
	"notes/internal/repositories"
	"notes/internal/services"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

var logger *slog.Logger

type application struct {
	logger   *slog.Logger
	wg       sync.WaitGroup
	confs    *configs.Config
	handlers *handlers.Handlers
}

func Init() {
	logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
	err := godotenv.Load(".env")
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

func (app *application) serveHttp() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.confs.HTTP_PORT),
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	shutDownErrChan := make(chan error)
	app.gracefulShutdown(srv, shutDownErrChan)
	app.logger.Info("starting server ...", slog.Group("server", "addr", srv.Addr, "environment", app.confs.ENV))
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	app.logger.Info("server stopped ...", slog.Group("server", "addr", srv.Addr))
	app.wg.Wait()
	return nil
}

func (app *application) gracefulShutdown(srv *http.Server, shutdownErrChan chan error) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, syscall.SIGTERM, syscall.SIGINT)
		<-quitChannel
		defer cancel()
		shutdownErrChan <- srv.Shutdown(ctx)
	}()
}

func run(logger *slog.Logger) error {
	conf := configs.New()

	db, err := db.New(logger, conf.DB_URI, conf.DB_NAME)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}

	httpErrs := handlers.NewHttpErrors(logger)

	categoryRepo := repositories.NewCategoryRepository(db, conf)
	noteRepo := repositories.NewNoteRepository(db, conf, categoryRepo)
	userRepo := repositories.NewUserRepository(db, conf)
	categoryService := services.NewCategoryService(categoryRepo)
	noteService := services.NewNoteService(noteRepo, categoryService)
	userService := services.NewUserService(userRepo, noteService)
	noteHandler := handlers.NewNoteHandler(noteService, httpErrs)
	categoryHandler := handlers.NewCategoryHandler(categoryService, httpErrs)
	userHandler := handlers.NewUserHandler(userService, httpErrs)

	hdls := handlers.New(noteHandler, categoryHandler, userHandler, logger, httpErrs)

	app := &application{
		logger:   logger,
		confs:    conf,
		handlers: hdls,
	}

	return app.serveHttp()
}
