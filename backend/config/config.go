package config

import (
	"Pinspire/backend/database"
	"Pinspire/backend/models"

	"gorm.io/gorm"
)

type Config struct {
	DB *gorm.DB
}

func Load() *Config {
	db := database.GetDB()

	db.AutoMigrate(
		&models.User{},
		&models.Location{},
		&models.Board{},
		&models.Pin{},
		&models.Product{},
		&models.UserProductLink{},
		&models.CountryMapping{},
	)

	return &Config{DB: db}
}

// func SomeFunction() {
// 	db := database.GetDB()
// 	// Use db as needed
// }
