package database

import (
	model "LoganXav/sori/app/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func MigrateDatabase() error {
    if err := DB.AutoMigrate(&model.User{}); err != nil {
        return fmt.Errorf("Cannot migrate table User: %w", err)
    }
    if err := DB.AutoMigrate(&model.UserToken{}); err != nil {
        return fmt.Errorf("Cannot migrate table UserToken: %w", err)
    }
    if err := DB.AutoMigrate(&model.File{}); err != nil {
        return fmt.Errorf("Cannot migrate table File: %w", err)
    }

    if err := SeedDatabase(); err != nil {
        return fmt.Errorf("Failed to seed database: %w", err)
    }

    return nil
}

func SeedDatabase() error {
    if err := DB.First(&model.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
        DB.Create(&model.User{
            Username:    "loganxav",
            Email:       "loganxav@email.com",
            Password:    "logan123",
            ValidatedAt: time.Date(2023, 1, 1, 10, 10, 10, 0, time.UTC),
        })
    }
    return nil
}