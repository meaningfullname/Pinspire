package service

import (
	"Pinspire/backend/models"
	"context"
	"gorm.io/gorm"
)

type CountryService struct {
	db *gorm.DB
}

func NewCountryService(db *gorm.DB) *CountryService {
	return &CountryService{db: db}
}

func (s *CountryService) SetUserCountry(ctx context.Context, code string) error {
	userID := ctx.Value("userID").(string)
	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("country_code", code).Error
}
