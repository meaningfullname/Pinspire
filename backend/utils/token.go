package utils

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(15 * 24 * time.Hour).Unix(), // 15 days expiration
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SEC")))
}

func SetTokenCookie(c *gin.Context, token string) {
	// Set the cookie (adjust secure flags as needed)
	c.SetCookie("token", token, 15*24*60*60, "/", "", false, true)
}
