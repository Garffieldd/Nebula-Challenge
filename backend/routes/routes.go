package routes

import (
	"net/http"

	"github.com/Nebula-Challenge/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	/*
		Defining the default route message
	*/
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})

	})

	/*
		Defining the route group for domain information
	*/
	router.GET("/domains-info", handlers.GetDomainsInformation)
	router.GET("/domains-info/:id", handlers.GetDomainsInformationByID)
	router.POST("/create-domain-info", handlers.PostDomainInformation)

}
