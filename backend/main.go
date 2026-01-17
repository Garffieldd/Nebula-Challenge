package main

import (
	"fmt"
	"log"

	"github.com/Nebula-Challenge/config"
	"github.com/Nebula-Challenge/routes"
	"github.com/gin-gonic/gin"
)

/*
init function to establish MongoDB connection before the main function runs
*/
func init() {
	if err := config.ConnectToMongo(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v ", err)
	}

	fmt.Println("Connected to MongoDB successfully")
}

/*
main is the entry point of the application, setting up the Gin router and routes
*/
func main() {

	defer config.CloseMongoConnection()
	router := gin.Default()

	routes.SetupRoutes(router)

	router.Run("localhost:8080")
}
