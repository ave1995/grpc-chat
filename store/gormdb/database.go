package gormdb

import (
	"context"
	"fmt"

	"github.com/ave1995/grpc-chat/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormConnection(ctx context.Context, config config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open: %w", err)
	}

	err = db.WithContext(ctx).AutoMigrate(&message{}, &outboxEvent{})
	if err != nil {
		return nil, fmt.Errorf("db.AutoMigrate: %w", err)
	}

	return db, nil
}
