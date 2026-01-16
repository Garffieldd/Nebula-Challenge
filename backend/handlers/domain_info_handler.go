package handlers

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
GetDomainsInformation handles the GET request to retrieve domain information
*/
func GetDomainsInformation(c *gin.Context) { // el parametro es el contexto de Gin dado por un puntero para capturar la peticion del cliente
	c.IndentedJSON(http.StatusOK, DomainsInfo)
}

/*
	postDomainInformation handles the POST request to add a new domain information register
*/

func PostDomainInformation(c *gin.Context) {
	var newDomainInfo DomainInformation

	if err := c.BindJSON(&newDomainInfo); err != nil { //referencia en memoria
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar el JSON: " + fmt.Sprint(err)})
		return
	}

	if newDomainInfo.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se ingreso un ID para el nuevo dominio analizado"})
		return
	}

	if newDomainInfo.Url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se ingreso una URL para el nuevo dominio analizado"})
		return
	}

	if newDomainInfo.Calification == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se ingreso una calificacion para el nuevo dominio analizado"})
		return
	}

	if newDomainInfo.Problems == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se ingresaron problemas para el nuevo dominio analizado"})
		return
	}

	if newDomainInfo.Details == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se ingresaron detalles para el nuevo dominio analizado"})
		return
	}

	DomainsInfo = append(DomainsInfo, newDomainInfo)

	c.IndentedJSON(http.StatusCreated, gin.H{"data": DomainsInfo, "message": "Domain information created successfully"})
}

func GetDomainsInformationByID(c *gin.Context) {
	id := c.Param("id")
	for _, domain_info := range DomainsInfo {
		if domain_info.ID == id {
			c.IndentedJSON(http.StatusOK, domain_info)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "domain information not found"}) //H es un mapa para enviar mensajes en JSON
}
