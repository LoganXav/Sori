package database

import (
	model "LoganXav/sori/app/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func MigrateDatabase() error{
	var err error

	if err = DB.AutoMigrate(&model.User{}); err != nil && DB.Migrator().HasTable(&model.User{}) {
		if err := DB.First(&model.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			DB.Create(&model.User{
				Username:    "loganxav",
				Email:       "loganxav@email.com",
				Password:    "logan123",
				ValidatedAt: time.Date(2023, 1, 1, 10, 10, 10, 0, time.UTC),
			})
		
		}
	}
	if err = DB.AutoMigrate(&model.UserToken{}); err != nil {
		return fmt.Errorf("Cannot migrate table UserToken")
	}

	if err = DB.AutoMigrate(&model.File{}); err != nil {
		return fmt.Errorf("Cannot migrate table File")
	}
	
	return nil
}