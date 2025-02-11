package main

import (
	"log"
	"net/http"

	"Pinspire/backend/database"
	"Pinspire/backend/models"
	"Pinspire/backend/service"

	"github.com/gorilla/mux"
)

func main() {
	// Get database connection
	db := database.GetDB()

	// Enable uuid-ossp extension
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatalf("Failed to create uuid extension: %v", err)
	}

	// Автоматическая миграция базы данных в правильном порядке
	err := db.AutoMigrate(
		&models.Location{},        // First, as it's referenced by User
		&models.User{},            // Second, as it's referenced by Board and Pin
		&models.Product{},         // Third, as it's referenced by Pin
		&models.Board{},           // Fourth, depends on User
		&models.Pin{},             // Fifth, depends on User, Board, and Product
		&models.UserProductLink{}, // Last, depends on User and Product
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Инициализация сервисов
	userService := service.NewUserService(db)
	fallbackService := service.NewFallbackService(db)

	// Создаем новый маршрутизатор mux
	r := mux.NewRouter()

	// Регистрируем маршруты для пользователей
	r.HandleFunc("/users", userService.CreateUser).Methods("POST")
	r.HandleFunc("/users", userService.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", userService.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", userService.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", userService.DeleteUser).Methods("DELETE")

	// Регистрируем fallback-сервис для отображения ссылок
	// Маршрут /display/ будет обрабатывать запросы вида:
	// GET /display/{productID}?location_id={locationID}
	r.HandleFunc("/display/", fallbackService.LinkDisplayHandler).Methods("GET")

	// Запуск сервера
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
