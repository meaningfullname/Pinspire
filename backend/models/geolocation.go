package models

import "gorm.io/gorm"

// Geolocation represents the data fetched from the geolocation API.
type Geolocation struct {
	gorm.Model
	Country   string `json:"country"`
	Region    string `json:"region"`
	City      string `json:"city"`
	IPAddress string `json:"ip_address"`
}
