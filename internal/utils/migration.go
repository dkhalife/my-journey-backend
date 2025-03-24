package utils

import (
	"dkhalife.com/journey/internal/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		models.Continent{},
		models.Country{},
		models.State{},
		models.City{},
		models.Address{},
		models.Position{},
		models.Tracker{},
		models.TrackingLocation{},
		models.User{},
		models.UserPasswordReset{},
		models.AppToken{},
	)
}
