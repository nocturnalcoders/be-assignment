package controllers

import (
	"concrete/configs"
	"concrete/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetPaymentAccounts retrieves payment accounts for a user
func GetPaymentAccounts(c *gin.Context) {
	// Get the user ID from the request params
	userID := c.Param("userId")

	db := configs.ConnectDB()
	// Get the payment accounts collection
	collection := db.Collection("payment_accounts")

	// Define a filter to find payment accounts for the given user ID
	filter := bson.M{"userId": userID}

	// Find payment accounts that match the filter
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payment accounts"})
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate over the cursor and collect payment accounts
	var paymentAccounts []models.PaymentAccount
	for cursor.Next(context.Background()) {
		var account models.PaymentAccount
		if err := cursor.Decode(&account); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode payment account"})
			return
		}
		paymentAccounts = append(paymentAccounts, account)
	}

	// Check if any payment accounts were found
	if len(paymentAccounts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No payment accounts found for the user"})
		return
	}

	// Return the payment accounts
	c.JSON(http.StatusOK, paymentAccounts)
}
