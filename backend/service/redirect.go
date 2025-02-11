package service

import (
	_ "fmt"
	"log"
	"net/http"
	"strings"

	"Pinspire/backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Redirect_service struct {
	db *gorm.DB
}

func (d *Redirect_service) RedirectHandler(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "Неверный формат URL", http.StatusBadRequest)
		return
	}
	productIDStr := parts[2]

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		http.Error(w, "Неверный product ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := d.db.First(&product, "product_id = ?", productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Продукт не найден", http.StatusNotFound)
			return
		}
		log.Printf("Ошибка при поиске продукта: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	locationIDStr := r.URL.Query().Get("location_id")
	if locationIDStr == "" {
		http.Error(w, "Параметр location_id обязателен", http.StatusBadRequest)
		return
	}
	locationID, err := uuid.Parse(locationIDStr)
	if err != nil {
		http.Error(w, "Неверный location_id", http.StatusBadRequest)
		return
	}

	var location models.Location
	if err := d.db.First(&location, "location_id = ?", locationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Локация не найдена", http.StatusNotFound)
			return
		}
		log.Printf("Ошибка при поиске локации: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	country := location.Country
	log.Printf("Используем страну из локации: %s", country)

	redirectURL := product.DefaultLink
	if country != "" && product.LocalizedLink != "" {
		redirectURL = product.LocalizedLink
	}

	log.Printf("Редирект на URL: %s", redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
