package utils

import (
	"dkhalife.com/journey/internal/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		models.User{},
	)
}
