// routes/routes.go
package routes

import (
	"concrete/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

// @title 	Tag Service API
// @version	1.0
// @description A Tag service API in Go using Gin framework

// @host 	localhost:8888
// @BasePath /api

func SetupRouter(db *mongo.Database) *gin.Engine {
	// var hmacSecret = "ovdi4ArU4JDlG9vrrxyOcS8Aa3sf690hhWIM4YkZr0awxZtI0GzeUbCrjW5AhtNbOoZEoTZ0M0rTIW8KRKF6Xw=="

	corsConfig := cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
	}

	router := gin.Default()
	url := ginSwagger.URL("http://localhost:8888/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.Use(cors.New(corsConfig))

	// Handle account creation
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.POST("/transactions", controllers.CreateTransaction)
	router.GET("/payment-accounts/:userId", controllers.GetPaymentAccounts)
	router.GET("/payment-history/:userId", controllers.GetPaymentHistory)

	return router
}
