package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"Pinspire/backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FallbackService struct {
	db *gorm.DB
}

func NewFallbackService(db *gorm.DB) *FallbackService {
	return &FallbackService{db: db}
}

func (fs *FallbackService) LinkDisplayHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем productID из URL (ожидаемый формат: /display/{productID})
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "Неверный формат URL", http.StatusBadRequest)
		return
	}
	productIDStr := parts[2]

	// Проверяем корректность формата productID (UUID)
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		http.Error(w, "Неверный product ID", http.StatusBadRequest)
		return
	}

	// Загружаем продукт из базы данных
	var product models.Product
	if err := fs.db.First(&product, "product_id = ?", productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Продукт не найден", http.StatusNotFound)
			return
		}
		log.Printf("Ошибка при поиске продукта: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Если передан параметр location_id, пытаемся загрузить локацию и извлечь страну
	var country string
	locationIDStr := r.URL.Query().Get("location_id")
	if locationIDStr != "" {
		if locationID, err := uuid.Parse(locationIDStr); err == nil {
			var location models.Location
			if err := fs.db.First(&location, "location_id = ?", locationID).Error; err == nil {
				country = location.Country
			} else {
				log.Printf("Локация с id %s не найдена: %v", locationIDStr, err)
			}
		} else {
			log.Printf("Неверный формат location_id: %v", err)
		}
	}

	// Выбираем, какую ссылку отображать:
	// Если страна определена и в продукте задана локализованная ссылка, используем её;
	// иначе, если задан дефолтный линк — используем его.
	var displayedLink string
	if country != "" && product.LocalizedLink != "" {
		displayedLink = product.LocalizedLink
	} else if product.DefaultLink != "" {
		displayedLink = product.DefaultLink
	}

	// Формируем ответ как JSON (используем анонимную map)
	response := map[string]interface{}{
		"product_id":   productID.String(),
		"product_name": product.Name,
		"country":      country,
	}
	if displayedLink == "" {
		response["fallback_message"] = "Ссылка недоступна для данного продукта. Пожалуйста, обратитесь в службу поддержки или попробуйте позже."
	} else {
		response["displayed_link"] = displayedLink
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка при кодировании ответа: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}
