package utils

import (
	"dkhalife.com/journey/internal/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		models.User{},
	); err != nil {
		return err
	}

	return nil
}
