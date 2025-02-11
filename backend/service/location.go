// service/handler.go
package service

import (
	"Pinspire/backend/models"
	_ "context"
	"encoding/json"
	"net/http"

	_ "github.com/gorilla/mux"
	"gorm.io/gorm"
)

type CountryRequest struct {
	CountryCode string `json:"country_code"`
}

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) SetUserCountry(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Парсинг запроса
	var req CountryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Валидация кода страны
	var exists bool
	if err := h.DB.Model(&models.CountryMapping{}).
		Select("count(*) > 0").
		Where("country_code = ?", req.CountryCode).
		Find(&exists).Error; err != nil || !exists {
		http.Error(w, "Invalid country code", http.StatusBadRequest)
		return
	}

	// Сохранение в БД в транзакции
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&models.User{}).
			Where("user_id = ?", userID).
			Update("manual_country", req.CountryCode).Error
	})

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
