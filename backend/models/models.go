package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	LocationID uuid.UUID `gorm:"type:uuid" json:"location_id"`
	LastActive time.Time `json:"last_active"`
}

type Location struct {
	ID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	Country string    `json:"country"`
	Region  string    `json:"region"`
	City    string    `json:"city"`
	IP      string    `json:"ip"`
}

type Board struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid" json:"user_id"`
	BoardName   string    `json:"board_name"`
	Description string    `json:"board_description"`
	IsPublic    bool      `json:"is_public"`
}

type Pin struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid" json:"user_id"`
	ProductID   uuid.UUID `gorm:"type:uuid" json:"product_id"`
	BoardID     uuid.UUID `gorm:"type:uuid" json:"board_id"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
}

type Product struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	Name          string    `json:"name"`
	DefaultLink   string    `json:"default_link"`
	LocalizedLink string    `json:"localized_link"`
	Price         float64   `json:"price"`
}

type UserProductLink struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid" json:"user_id"`
	ProductID uuid.UUID `gorm:"type:uuid" json:"product_id"`
	ViewedAt  time.Time `json:"viewed_at"`
}
