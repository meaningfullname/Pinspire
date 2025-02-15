package main

import (
	"log"
	"os"

	"Pinspire/backend/database"
	"Pinspire/backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MongoDB
	dbClient, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Database connection error: ", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Make the database client available in the context
	router.Use(func(c *gin.Context) {
		c.Set("db", dbClient)
		c.Next()
	})

	// Register API routes
	routes.RegisterRoutes(router)

	// Serve static frontend files (assumes you built your frontend into ./frontend/dist)
	router.Static("/", "./frontend/dist")
	router.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	router.Run(":" + port)
}
