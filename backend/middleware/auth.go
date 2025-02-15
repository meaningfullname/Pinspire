package middleware

import (
	"context"
	"net/http"
	"os"
	"time"

	"Pinspire/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Please Login"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SEC")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, err := primitive.ObjectIDFromHex(claims["id"].(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid user id in token"})
			c.Abort()
			return
		}

		// Get the database client from context
		dbClient, exists := c.Get("db")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database connection missing"})
			c.Abort()
			return
		}
		client := dbClient.(*mongo.Client)
		collection := client.Database("pinterest").Collection("users")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		err = collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
			c.Abort()
			return
		}

		// Set the user in context for later handlers
		c.Set("user", user)
		c.Next()
	}
}
