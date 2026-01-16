package main

import (
	"github.com/Nebula-Challenge/routes"
	"github.com/gin-gonic/gin"
)

/*
main is the entry point of the application, setting up the Gin router and routes
*/
func main() {

	router := gin.Default()

	routes.SetupRoutes(router)

	router.Run("localhost:8080")
}
