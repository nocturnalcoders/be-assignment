package configs

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB establishes a connection to the MongoDB database
func ConnectDB() *mongo.Database {
	// Set up MongoDB connection options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection to the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database Successfully Connected")
	}

	// Return the database instance
	return client.Database("mydb")
}
