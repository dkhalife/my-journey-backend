package models

import "time"

type User struct {
	ID          int       `json:"id" gorm:"primary_key;not null"`
	DisplayName string    `json:"display_name" gorm:"column:display_name;not null"`
	Email       string    `json:"email" gorm:"column:email;unique;not null"`
	Password    string    `json:"-" gorm:"column:password;not null"`
	CreatedAt   time.Time `json:"-" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"-" gorm:"column:updated_at;default:NULL;autoUpdateTime"`
	Disabled    bool      `json:"-" gorm:"column:disabled;default:false"`
}
