package models

// Redirection defines the mapping between a country and a URL.
type Redirection struct {
	CountryCode string
	URL         string
}

// DefaultLinks is the predefined map of links for redirection.
var DefaultLinks = map[string]string{
	"KZ":      "https://www.kaspi.kz",
	"US":      "https://www.amazon.com",
	"RU":      "https://www.wildberries.ru",
	"DEFAULT": "https://www.global.com",
}
