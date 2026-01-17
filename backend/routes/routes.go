package routes

import (
	"net/http"

	"github.com/Nebula-Challenge/handlers"
	"github.com/gin-gonic/gin"
)

/*
SetupRoutes configures the routes for the Gin router
Args:

	router *gin.Engine: The Gin router instance to set up routes on

Returns:

	None
*/
func SetupRoutes(router *gin.Engine, handler *handlers.Handler) {
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
	// handler := &handlers.Handler{DB: config.GetMongoClient()}
	router.GET("/domains-info", handler.GetDomainsInformation)
	router.GET("/domains-info/:id", handler.GetDomainsInformationByID)
	router.POST("/create-domain-info", handler.PostDomainInformation)
	router.POST("/domains-info/aggregate", handler.AggregateDomainInformation)
	router.DELETE("/domains-info/:id", handler.DeleteDomainById)

	//SSL Labs TLS scan routes
	router.POST("/start-scan", handler.StartScan)
	router.GET("/scan-status/:scanRequestID", handler.GetScanStatus)

}
