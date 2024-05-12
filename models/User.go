package models

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID       int    `json:"id" bson:"_id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// DTO
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

// JwtClaims represents the JWT claims
type JwtClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Valid implements jwt.Claims.
// Subtle: this method shadows the method (StandardClaims).Valid of JwtClaims.StandardClaims.
func (j JwtClaims) Valid() error {
	panic("unimplemented")
}

func NewUser(username, password, email string) *User {
	// In a real application, you would generate the ID using a more sophisticated method
	// For simplicity, we'll just increment a global variable
	userID := getNextUserID()

	return &User{
		ID:       userID,
		Username: username,
		Password: password,
		Email:    email,
	}
}

// Dummy implementation of auto-incrementing ID generation
var nextUserID int = 1

func getNextUserID() int {
	id := nextUserID
	nextUserID++
	return id
}
