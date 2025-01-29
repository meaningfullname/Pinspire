package utils

import (
	"Pinspire/backend/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// FetchGeolocation fetches the user's geolocation based on their IP.
func FetchGeolocation(ip string) (*models.Geolocation, error) {
	apiURL := fmt.Sprintf("https://ipinfo.io/%s/json", ip)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var geo models.Geolocation
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
		return nil, err
	}
	return &geo, nil
}
