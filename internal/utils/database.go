package utils

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("journey.db"), &gorm.Config{})
}
