package main

import (
	"concrete/routes"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	db *mongo.Database
)

func main() {
	r := routes.SetupRouter(db)

	// Start the Gin server
	err := r.Run("0.0.0.0:8888")
	if err != nil {
		fmt.Println("Failed to start Gin server:", err)
	}

}
