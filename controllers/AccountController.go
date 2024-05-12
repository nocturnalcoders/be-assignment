package controllers

import (
	"concrete/configs"
	lib "concrete/libs"
	"concrete/models"
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateTags		godoc
// @Summary			Create tags
// @Description		Save tags data in Db.
// @Body			tags body model.RegisterRequest true "Create tags"
// @Produce			application/json
// @Tags			tags
// @Router			/tags [post]
func Register(c *gin.Context) {
	var requestBody models.RegisterRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := lib.HashPassword(requestBody.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Connect to MongoDB
	db := configs.ConnectDB()

	// Insert user data into MongoDB
	userCollection := db.Collection("users")
	user := models.User{
		Email:    requestBody.Email,
		Username: requestBody.Username,
		Password: hashedPassword,
	}
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// You can continue with Supabase integration here if needed

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Login handles user login
func Login(c *gin.Context) {
	var requestBody models.LoginRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate user
	user, err := authenticateUser(requestBody.Email, requestBody.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	tokenString, err := GenerateJWTToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return token in response
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Authenticate user based on email and password
func authenticateUser(email, password string) (*models.User, error) {
	// Connect to your database
	db := configs.ConnectDB()

	// Query the database to find the user by email
	userCollection := db.Collection("users")
	filter := bson.M{"email": email}
	var user models.User
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err // User not found or database error
	}

	// Verify the password
	if err := lib.VerifyPassword(user.Password, password); err != nil {
		return nil, err // Password doesn't match
	}

	// Authentication successful, return the user
	return &user, nil
}

// GenerateJWTToken generates a JWT token with the given email
func GenerateJWTToken(email string) (string, error) {
	// Define JWT claims
	claims := models.JwtClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte("C0NC3R3TE"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
