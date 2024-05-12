package controllers

import (
	"concrete/configs"
	"concrete/models"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetTransactions handles retrieving transactions for an account
func GetTransactions(c *gin.Context, db *mongo.Database) {
	// Get the account ID from the request params
	accountID := c.Param("accountId")

	// Get the transactions collection
	collection := db.Collection("transactions")

	// Define a filter to find transactions for the given account ID
	filter := bson.M{"accountId": accountID}

	// Find transactions that match the filter
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate over the cursor and collect transactions
	var transactions []models.Transaction
	for cursor.Next(context.Background()) {
		var transaction models.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode transaction"})
			return
		}
		transactions = append(transactions, transaction)
	}

	// Check if any transactions were found
	if len(transactions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No transactions found for the account"})
		return
	}

	// Return the transactions
	c.JSON(http.StatusOK, transactions)
}

// CreateTransaction handles creating a new transaction
func CreateTransaction(c *gin.Context) {
	var requestBody models.Transaction
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Connect to the database
	db := configs.ConnectDB()
	ctx := context.TODO()

	// Start a MongoDB session
	session, err := db.Client().StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start session"})
		return
	}
	defer session.EndSession(ctx)

	// Start a transaction
	err = session.StartTransaction()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Get the transactions collection
	transactionCollection := db.Collection("transactions")
	// Get the payment accounts collection
	accountCollection := db.Collection("payment_accounts")
	// Get the payment history collection
	historyCollection := db.Collection("payment_history")

	// Insert the new transaction into the database
	_, err = transactionCollection.InsertOne(ctx, requestBody)
	if err != nil {
		session.AbortTransaction(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Retrieve the payment account from the database
	var paymentAccount models.PaymentAccount
	err = accountCollection.FindOne(ctx, bson.M{"_id": requestBody.AccountID}).Decode(&paymentAccount)
	if err != nil {
		// Log the error for debugging purposes
		fmt.Println("Error retrieving payment account:", err)

		session.AbortTransaction(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payment account"})
		return
	}
	// Update the balance of the associated payment account
	newBalance := paymentAccount.Balance + requestBody.Amount
	_, err = accountCollection.UpdateOne(ctx, bson.M{"_id": requestBody.AccountID}, bson.M{"$set": bson.M{"balance": newBalance}})
	if err != nil {
		session.AbortTransaction(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment account balance"})
		return
	}

	// Create a new payment history record
	paymentHistory := models.PaymentHistory{
		AccountID:   requestBody.AccountID,
		Amount:      requestBody.Amount,
		Currency:    requestBody.Currency,
		Description: requestBody.Description,
		Timestamp:   requestBody.Timestamp,
	}
	_, err = historyCollection.InsertOne(ctx, paymentHistory)
	if err != nil {
		session.AbortTransaction(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment history"})
		return
	}

	// Commit the transaction
	err = session.CommitTransaction(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully"})
}
