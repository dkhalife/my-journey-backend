package models

import "time"

type Tracker struct {
	ID        int       `json:"id" gorm:"primary_key;not null"`
	Name      string    `json:"name" gorm:"column:name;not null"`
	UserID    int       `json:"user_id" gorm:"column:user_id;not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time `json:"-" gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
}

type TrackingMetadata struct {
	Accuracy         float64 `json:"accuracy" gorm:"column:accuracy;not null"`
	VerticalAccuracy float64 `json:"vertical_accuracy" gorm:"column:vertical_accuracy"`
	Course           float64 `json:"course" gorm:"column:course"`
	Velocity         float64 `json:"velocity" gorm:"column:velocity"`
	BatteryLevel     int     `json:"battery_level" gorm:"column:battery_level"`
	SSID             string  `json:"ssid" gorm:"column:ssid"`
	BSSID            string  `json:"bssid" gorm:"column:bssid"`
}

type TrackingLocation struct {
	ID         int              `json:"id" gorm:"primary_key;not null"`
	TimeStamp  time.Time        `json:"timestamp" gorm:"column:timestamp;not null"`
	UserID     int              `json:"user_id" gorm:"column:user_id;not null"`
	User       User             `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TrackerID  int              `json:"tracker_id" gorm:"column:tracker_id"`
	Tracker    Tracker          `json:"tracker" gorm:"foreignKey:TrackerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PositionID int              `json:"position_id" gorm:"column:position_id;not null"`
	Position   Position         `json:"position" gorm:"foreignKey:PositionID;references:ID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	Metadata   TrackingMetadata `json:"metadata" gorm:"embedded;embeddedPrefix:metadata_"`
	CreatedAt  time.Time        `json:"-" gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt  time.Time        `json:"-" gorm:"column:updated_at;default:NULL;autoUpdateTime;not null"`
}
