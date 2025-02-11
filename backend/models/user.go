package models

import (
	"time"

	"github.com/google/uuid"
)

// User соответствует таблице users (столбец: user_id)
type User struct {
	ID         uuid.UUID `gorm:"column:user_id;type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`
	Email      string    `gorm:"column:email;unique" json:"email"`
	Password   string    `gorm:"column:password" json:"-"`
	LocationID uuid.UUID `gorm:"column:location_id;type:uuid" json:"location_id"`
	LastActive time.Time `gorm:"column:last_active" json:"last_active"`
	Location   Location  `gorm:"foreignKey:LocationID;references:ID" json:"location,omitempty"`
	Boards     []Board   `gorm:"foreignKey:UserID;references:ID" json:"boards,omitempty"`
	Pins       []Pin     `gorm:"foreignKey:UserID;references:ID" json:"pins,omitempty"`
	Products   []Product `gorm:"many2many:user_product_links;" json:"products,omitempty"`
}

// Задаём имя таблицы для User
func (User) TableName() string {
	return "users"
}

// Location соответствует таблице locations (столбец: location_id)
type Location struct {
	ID      uuid.UUID `gorm:"column:location_id;type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Country string    `gorm:"column:country" json:"country"`
	Region  string    `gorm:"column:region" json:"region"`
	City    string    `gorm:"column:city" json:"city"`
	IP      string    `gorm:"column:ip_address" json:"ip"` // ip_address в SQL-скрипте
}

func (Location) TableName() string {
	return "locations"
}

// Board соответствует таблице boards (столбец: board_id)
type Board struct {
	ID          uuid.UUID `gorm:"column:board_id;type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid" json:"user_id"`
	BoardName   string    `gorm:"column:board_name" json:"board_name"`
	Description string    `gorm:"column:board_description" json:"board_description"`
	IsPublic    bool      `gorm:"column:is_public" json:"is_public"`
	// Если нужно, можно добавить связь с User (но избегайте циклических зависимостей в JSON)
	User User  `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Pins []Pin `gorm:"foreignKey:BoardID;references:ID" json:"pins,omitempty"`
}

func (Board) TableName() string {
	return "boards"
}

// Pin соответствует таблице pins (столбец: pin_id)
type Pin struct {
	ID          uuid.UUID `gorm:"column:pin_id;type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid" json:"user_id"`
	ProductID   uuid.UUID `gorm:"column:product_id;type:uuid" json:"product_id"`
	BoardID     uuid.UUID `gorm:"column:board_id;type:uuid" json:"board_id"`
	ImageURL    string    `gorm:"column:image_url" json:"image_url"`
	Description string    `gorm:"column:description" json:"description"`
	Tags        []string  `gorm:"column:tags;type:text[]" json:"tags"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	// Связи (опционально, если не требуются в JSON, их можно убрать или настроить json:"-")
	User    User    `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Board   Board   `gorm:"foreignKey:BoardID;references:ID" json:"-"`
	Product Product `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
}

func (Pin) TableName() string {
	return "pins"
}

// Product соответствует таблице products (столбец: product_id)
type Product struct {
	ID            uuid.UUID `gorm:"column:product_id;type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name          string    `gorm:"column:product_name" json:"name"`
	DefaultLink   string    `gorm:"column:default_link" json:"default_link"`
	LocalizedLink string    `gorm:"column:localized_link" json:"localized_link"`
	Price         float64   `gorm:"column:price" json:"price"`
}

func (Product) TableName() string {
	return "products"
}

// UserProductLink соответствует таблице user_product_links (столбец: user_product_link_id)
type UserProductLink struct {
	ID        uuid.UUID `gorm:"column:user_product_link_id;type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid" json:"user_id"`
	ProductID uuid.UUID `gorm:"column:product_id;type:uuid" json:"product_id"`
	ViewedAt  time.Time `gorm:"column:viewed_at" json:"viewed_at"`
}

func (UserProductLink) TableName() string {
	return "user_product_links"
}

// CountryMapping остаётся без изменений, если схема соответствует
type CountryMapping struct {
	TimezonePattern string `gorm:"column:timezone_pattern" json:"timezone_pattern"`
	LanguagePattern string `gorm:"column:language_pattern" json:"language_pattern"`
	CurrencyPattern string `gorm:"column:currency_pattern" json:"currency_pattern"`
	CountryCode     string `gorm:"column:country_code;primaryKey" json:"country_code"`
}

func (CountryMapping) TableName() string {
	return "country_mappings" // если у вас таблица с таким именем
}

type LinkDisplayResponse struct {
	ProductID       string `json:"product_id"`
	ProductName     string `json:"product_name"`
	DisplayedLink   string `json:"displayed_link,omitempty"`
	FallbackMessage string `json:"fallback_message,omitempty"`
	Country         string `json:"country,omitempty"`
}
