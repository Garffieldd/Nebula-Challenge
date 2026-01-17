package main

import (
	"fmt"
	"log"
	"github.com/Nebula-Challenge/config"
	"github.com/Nebula-Challenge/handlers"
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

	// Good practice to close the MongoDB connection when the application exits
	defer config.CloseMongoConnection()

	// Setting up the Gin router and routes
	dbConfig := config.GetMongoClient()
	handler := handlers.NewHandler(dbConfig)
	router := gin.Default()
	routes.SetupRoutes(router, handler)

	router.Run("localhost:8080")
}
