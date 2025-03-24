package models

type Continent struct {
	ID   int    `json:"id" gorm:"primary_key;not null"`
	Name string `json:"name" gorm:"column:name;not null"`
}

type Country struct {
	ID          int       `json:"id" gorm:"primary_key;not null"`
	Name        string    `json:"name" gorm:"column:name;not null"`
	Code        string    `json:"code" gorm:"column:code;not null"`
	ContinentID int       `json:"continent_id" gorm:"column:continent_id;not null"`
	Continent   Continent `json:"continent" gorm:"foreignKey:ContinentID;references:ID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
}

type State struct {
	ID        int     `json:"id" gorm:"primary_key;not null"`
	Name      string  `json:"name" gorm:"column:name;not null"`
	Code      string  `json:"code" gorm:"column:code;not null"`
	CountryID int     `json:"country_id" gorm:"column:country_id;not null"`
	Country   Country `json:"country" gorm:"foreignKey:CountryID;references:ID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
}

type City struct {
	ID      int    `json:"id" gorm:"primary_key;not null"`
	Name    string `json:"name" gorm:"column:name;not null"`
	StateID int    `json:"state_id" gorm:"column:state_id;not null"`
	State   State  `json:"state" gorm:"foreignKey:StateID;references:ID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
}

type Address struct {
	ID           int    `json:"id" gorm:"primary_key;not null"`
	StreetNumber string `json:"street_number" gorm:"column:street_number;not null"`
	StreetName   string `json:"street_name" gorm:"column:street_name;not null"`
	CityID       int    `json:"city_id" gorm:"column:city_id;not null"`
	City         City   `json:"city" gorm:"foreignKey:CityID;references:ID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
}

type Position struct {
	ID        int     `json:"id" gorm:"primary_key;not null"`
	Longitude float64 `json:"longitude" gorm:"column:longitude;not null"`
	Latitude  float64 `json:"latitude" gorm:"column:latitude;not null"`
	Altitude  float64 `json:"altitude" gorm:"column:altitude;not null"`
	AddressID int     `json:"address_id" gorm:"column:address_id;not null"`
	Address   Address `json:"address" gorm:"foreignKey:AddressID;references:ID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
}
