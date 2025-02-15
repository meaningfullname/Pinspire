package controllers

import (
	"context"
	"net/http"
	"time"

	"Pinspire/backend/models"
	"Pinspire/backend/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	client := c.MustGet("db").(*mongo.Client)
	collection := client.Database("pinterest").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user exists
	var existing models.User
	err := collection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Already have an account with this email"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		return
	}

	newUser := models.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  string(hashed),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating user"})
		return
	}
	newUser.ID = result.InsertedID.(primitive.ObjectID)

	token, err := utils.GenerateToken(newUser.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating token"})
		return
	}
	utils.SetTokenCookie(c, token)

	c.JSON(http.StatusCreated, gin.H{"user": newUser, "message": "User Registered"})
}

func LoginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	client := c.MustGet("db").(*mongo.Client)
	collection := client.Database("pinterest").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No user with this email"})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong password"})
		return
	}

	token, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating token"})
		return
	}
	utils.SetTokenCookie(c, token)
	c.JSON(http.StatusOK, gin.H{"user": user, "message": "Logged in"})
}

func MyProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UserProfile(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}

	client := c.MustGet("db").(*mongo.Client)
	collection := client.Database("pinterest").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}
	user.Password = "" // omit password
	c.JSON(http.StatusOK, user)
}

// FollowAndUnfollowUser and LogOutUser can be implemented in a similar fashion.
