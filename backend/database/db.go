package database

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
}

// GetDB returns a singleton instance of the database connection
func GetDB() *gorm.DB {
	once.Do(func() {
		config := getConfig()

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			config.Host,
			config.User,
			config.Password,
			config.DBName,
			config.Port,
			config.SSLMode,
		)

		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	})
	return db
}

func getConfig() Config {
	return Config{
		Host:     getEnv("DB_HOST", "localhost"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "tryu1234"),
		DBName:   getEnv("DB_NAME", "pinspire"),
		Port:     getEnv("DB_PORT", "5432"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
