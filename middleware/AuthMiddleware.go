package middleware

import (
	"concrete/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(hmacSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		splitToken := strings.Split(token, "Bearer ")
		if len(splitToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		token = splitToken[1]

		email, err := parseJWTToken(token, []byte(hmacSecret))

		if err != nil {
			log.Printf("Error parsing token: %s", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		log.Printf("Received request from %s", email)

		ctx := context.WithValue(c, "email", email)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func parseJWTToken(token string, hmacSecret []byte) (email string, err error) {
	// Parse the token and validate the signature
	t, err := jwt.ParseWithClaims(token, &models.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	// Check if the token is valid
	if err != nil {
		return "", fmt.Errorf("error validating token: %v", err)
	} else if claims, ok := t.Claims.(*models.JwtClaims); ok {
		return claims.Email, nil
	}

	return "", fmt.Errorf("error parsing token: %v", err)
}
