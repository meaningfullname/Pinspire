package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/your_username/my_project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func SetUserCollection(client *mongo.Client) {
	userCollection = client.Database("mydatabase").Collection("users")
}

// Create User
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user.ID = primitive.NewObjectID()
	user.LastActive = time.Now()

	_, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// Get Users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	cursor, err := userCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	for cursor.Next(context.TODO()) {
		var user models.User
		cursor.Decode(&user)
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}
