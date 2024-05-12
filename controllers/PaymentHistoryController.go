package controllers

import (
	"concrete/configs"
	"concrete/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetPaymentHistory retrieves payment history for a user
func GetPaymentHistory(c *gin.Context) {
	// Get the user ID from the request params
	userID := c.Param("userId")

	db := configs.ConnectDB()
	// Get the payment history collection
	collection := db.Collection("payment_history")

	// Define a filter to find payment history for the given user ID
	filter := bson.M{"userId": userID}

	// Find payment history that match the filter
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payment history"})
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate over the cursor and collect payment history
	var paymentHistory []models.PaymentHistory
	for cursor.Next(context.Background()) {
		var history models.PaymentHistory
		if err := cursor.Decode(&history); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode payment history"})
			return
		}
		paymentHistory = append(paymentHistory, history)
	}

	// Check if any payment history were found
	if len(paymentHistory) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No payment history found for the user"})
		return
	}

	// Return the payment history
	c.JSON(http.StatusOK, paymentHistory)
}
