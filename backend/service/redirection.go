package services

import (
	"Pinspire/backend/models"
	"errors"
)

// GetRedirectURL returns the appropriate URL based on the country code.
func GetRedirectURL(countryCode string) (string, error) {
	url, exists := models.DefaultLinks[countryCode]
	if !exists {
		url = models.DefaultLinks["DEFAULT"] // Fallback link
	}
	return url, nil
}

// Example service to simulate link redirection based on country
func RedirectByGeo(countryCode string) (string, error) {
	if countryCode == "" {
		return "", errors.New("invalid country code")
	}
	return GetRedirectURL(countryCode)
}
