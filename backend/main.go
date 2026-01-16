package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
DomainInformation represents the structure for domain data
*/
type DomainInformation struct {
	ID           string `json:"id"`
	Url          string `json:"url"`
	Calification string `json:"calification"`
	Problems     string `json:"problems"`
	Details      string `json:"details"`
}

/*
DomainsInfo is a sample list of domain information for testing purposes
*/
var DomainsInfo = []DomainInformation{
	{
		ID:           "1",
		Url:          "example.com",
		Calification: "A",
		Problems:     "None",
		Details:      "This is a sample domain.",
	},
	{
		ID:           "2",
		Url:          "test.com",
		Calification: "B",
		Problems:     "Minor issues",
		Details:      "This is another sample domain.",
	},
	{
		ID:           "3",
		Url:          "sample.org",
		Calification: "C",
		Problems:     "Major issues",
		Details:      "This is yet another sample domain.",
	},
}

/*
getDomainsInformation handles the GET request to retrieve domain information
*/
func getDomainsInformation(c *gin.Context) { // el parametro es el contexto de Gin dado por un puntero para capturar la peticion del cliente
	c.IndentedJSON(http.StatusOK, DomainsInfo)
}

/*
	postDomainInformation handles the POST request to add a new domain information register
*/

func postDomainInformation(c *gin.Context) {
	var newDomainInfo DomainInformation

	if err := c.BindJSON(&newDomainInfo); err != nil { //referencia en memoria
		fmt.Println("Hubo un erro almacenando el nuevo registro en memoria")
		return
	}

	DomainsInfo = append(DomainsInfo, newDomainInfo)

	c.IndentedJSON(http.StatusCreated, DomainsInfo)
}

func getDomainsInformationByID(c *gin.Context) {
	id := c.Param("id")
	for _, domain_info := range DomainsInfo {
		if domain_info.ID == id {
			c.IndentedJSON(http.StatusOK, domain_info)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "domain information not found"}) //H es un mapa para enviar mensajes en JSON
}

/*
main is the entry point of the application, setting up the Gin router and routes
*/
func main() {
	router := gin.Default()

	/*
	Defining the default route message
	*/
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	router.GET("/domains-info", getDomainsInformation)
	router.GET("/domains-info/:id", getDomainsInformationByID)
	router.POST("/create-domain-info", postDomainInformation)

	router.Run("localhost:8080")
}
