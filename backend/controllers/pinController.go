package controllers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"Pinspire/backend/models"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// Make sure to import "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

var cld, _ = cloudinary.NewFromParams(
	os.Getenv("Cloud_Name"),
	os.Getenv("Cloud_Api"),
	os.Getenv("Cloud_Secret"),
)

func CreatePin(c *gin.Context) {
	title := c.PostForm("title")
	pinText := c.PostForm("pin")

	// Get file from form (similar to multer.memoryStorage)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "File is required"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error opening file"})
		return
	}
	defer file.Close()

	// Read file into a buffer
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading file"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	uploadResult, err := cld.Upload.Upload(ctx, &buf, uploader.UploadParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cloudinary upload failed", "error": err.Error()})
		return
	}

	// Get current user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	currentUser := userInterface.(models.User)

	// Insert pin into MongoDB
	client := c.MustGet("db").(*mongo.Client)
	collection := client.Database("pinterest").Collection("pins")
	newPin := models.Pin{
		Title:   title,
		PinText: pinText,
		Owner:   currentUser.ID,
		Image: models.Image{
			ID:  uploadResult.PublicID,
			URL: uploadResult.SecureURL,
		},
		Comments:  []models.Comment{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	insertCtx, insertCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer insertCancel()
	result, err := collection.InsertOne(insertCtx, newPin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating pin"})
		return
	}
	newPin.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusOK, gin.H{"message": "Pin Created", "pin": newPin})
}

// Implement GetAllPins, GetSinglePin, CommentOnPin, DeleteComment, DeletePin, UpdatePin similarly.
