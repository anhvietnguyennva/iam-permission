package db

import (
	"time"

	"github.com/anhvietnguyennva/go-error/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"iam-permission/internal/pkg/config"
)

var db *gorm.DB

func InitDB() error {
	if db != nil {
		return nil
	}

	cfg := config.Instance().DB

	gormDB, err := gorm.Open(
		postgres.Open(cfg.DSN()),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(cfg.LogLevel)),
		},
	)
	if err != nil {
		return errors.NewInfraErrorDBConnect(err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return errors.NewInfraErrorDBConnect(err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnLifeTime) * time.Second)

	if err = sqlDB.Ping(); err != nil {
		return errors.NewInfraErrorDBConnect(err)
	}

	db = gormDB

	return nil
}

func Instance() *gorm.DB {
	return db
}
