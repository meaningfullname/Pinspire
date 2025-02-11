package service

import (
	"encoding/json"
	"net/http"
	"time"

	"Pinspire/backend/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

// Generic CRUD Helper Functions
func (s *UserService) createRecord(w http.ResponseWriter, r *http.Request, model interface{}) {
	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := s.db.Create(model).Error; err != nil {
		http.Error(w, "Failed to insert record", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Record created successfully"})
}

func (s *UserService) getRecords(w http.ResponseWriter, model interface{}) {
	if err := s.db.Find(model).Error; err != nil {
		http.Error(w, "Error fetching records", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(model)
}

func (s *UserService) getRecordByID(w http.ResponseWriter, r *http.Request, model interface{}, idKey string) {
	params := mux.Vars(r)
	if err := s.db.First(model, "id = ?", params[idKey]).Error; err != nil {
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(model)
}

func (s *UserService) updateRecord(w http.ResponseWriter, r *http.Request, model interface{}, idKey string) {
	params := mux.Vars(r)
	if err := s.db.First(model, "id = ?", params[idKey]).Error; err != nil {
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := s.db.Save(model).Error; err != nil {
		http.Error(w, "Failed to update record", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Record updated successfully"})
}

func (s *UserService) deleteRecord(w http.ResponseWriter, r *http.Request, model interface{}, idKey string) {
	params := mux.Vars(r)
	if err := s.db.Delete(model, "id = ?", params[idKey]).Error; err != nil {
		http.Error(w, "Failed to delete record", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Record deleted successfully"})
}

// Update the handler methods to use the service methods
func (s *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	user.ID = uuid.New()
	user.LastActive = time.Now()
	s.createRecord(w, r, &user)
}

func (s *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	s.getRecords(w, &users)
}

func (s *UserService) GetUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	s.getRecordByID(w, r, &user, "id")
}

func (s *UserService) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	s.updateRecord(w, r, &user, "id")
}

func (s *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	s.deleteRecord(w, r, &user, "id")
}

// Continue updating all other handlers similarly...

// NewUserService creates a new UserService instance
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}
