package db

import (
	"context"
	"log"
	"log/slog"
	"notes/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	defaultTimeout     = time.Second * 3
	maxOpenConnections = 10
	maxIdleConnections = 5
)

var newLogger = logger.New(
	log.New(log.Writer(), "", log.LstdFlags),
	logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Silent,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	},
)

func New(logger *slog.Logger, dsn, dbName string) (*gorm.DB, error) {

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		logger.Error("GormDB Open failed", "error: ", err)
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Error("SQL DB failed", "error: ", err)
		return nil, err
	}

	sqlDB.SetConnMaxLifetime(defaultTimeout)
	sqlDB.SetMaxOpenConns(maxOpenConnections)
	sqlDB.SetMaxIdleConns(maxIdleConnections)

	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Error("DB Ping Test Failed", "error: ", err)
		return nil, err
	}

	logger.Info("Running migrations...")
	if err := gormDB.AutoMigrate(&models.User{}, &models.Note{}, &models.Category{}); err != nil {
		logger.Error("AutoMigrate failed", "error", err)
		return nil, err
	}

	logger.Info("Migrations completed successfully.")
	logger.Info("Successfully connected to DB: " + dbName)

	return gormDB, nil
}
