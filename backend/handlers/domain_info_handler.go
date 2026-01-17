package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Nebula-Challenge/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Handler struct to hold the MongoDB client instance
*/
type Handler struct {
	DB *config.DatabaseConfig
}

/*
GetDomainsInformation handles the GET request to retrieve domain information, fetching data from MongoDB

Args:

	c *gin.Context: The Gin context for handling the request and response

Returns:

	None: Sends a JSON response with the list of domain information or an error message
*/
func (h *Handler) GetDomainsInformation(c *gin.Context) { // el parametro es el contexto de Gin dado por un puntero para capturar la peticion del cliente

	coll := h.DB.Client.Database(h.DB.DbName).Collection("domains_info")

	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obtaining MongoDB data: " + fmt.Sprint(err)})
		return
	}

	var domainsInfo []bson.M                                        // Resultados de la consulta de arriva
	if err = cursor.All(context.TODO(), &domainsInfo); err != nil { //Referencia a la varuable domainsInfo donde se almacenaran los datos obtenidos
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding MongoDB data: " + fmt.Sprint(err)})
		return
	}

	c.JSON(http.StatusOK, domainsInfo)
}

/*
	postDomainInformation handles the POST request to add a new domain information register
Args:

	c *gin.Context: The Gin context for handling the request and response
Returns:

	None: Sends a JSON response confirming insertion or an error message

*/

func (h *Handler) PostDomainInformation(c *gin.Context) {
	var pipeline interface{}                            //Donde se almacenara la informacion recibida
	if err := c.ShouldBindJSON(&pipeline); err != nil { //Sirve para descerializar una solicitud POST
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error decoding JSON: " + fmt.Sprint(err)})
		return
	}

	coll := h.DB.Client.Database(h.DB.DbName).Collection("domains_info")
	result, err := coll.InsertOne(context.TODO(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting document: " + fmt.Sprint(err)})
		return
	}

	insertedID := result.InsertedID

	c.JSON(http.StatusCreated, gin.H{"message": "Document inserted successfully", "inserted_id": insertedID})

}

/*
GetDomainsInformationByID handles the GET request to retrieve domain information by its ID, fetching data from MongoDB
Args:

	c *gin.Context: The Gin context for handling the request and response

Returns:

	None: Sends a JSON response with the domain information or an error message
*/
func (h *Handler) GetDomainsInformationByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr) //convierte el string a ObjectID (porque asi lo maneja mongo)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format: " + fmt.Sprint(err)})
		return
	}

	var domainInfo bson.M //mapa para almacenar el resultado

	coll := h.DB.Client.Database(h.DB.DbName).Collection("domains_info")
	err = coll.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&domainInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obtaining MongoDB data: " + fmt.Sprint(err)})
		return
	}
	c.JSON(http.StatusOK, domainInfo)
}

/*
AggregateDomainInformation handles the POST request to perform aggregation on domain information in MongoDB
A body example:
[

	  { "$match": {
	   "Calification": "A"
	   }
	   },
	   {
	   "$sort": {
	    "count": -1
	}

	}

]

This instruction returns a list of the domains with calification "A" sorted in descending order.

Args:

	c *gin.Context: The Gin context for handling the request and response

Returns:

	None: Sends a JSON response with the aggregation result or an error message
*/
func (h *Handler) AggregateDomainInformation(c *gin.Context) {

	var pipeline interface{}                            //Para que pueda contener cualquier tipo de valor
	if err := c.ShouldBindJSON(&pipeline); err != nil { //Sirve para descerializar una solicitud POST
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error decoding JSON: " + fmt.Sprint(err)})
		return
	}
	coll := h.DB.Client.Database(h.DB.DbName).Collection("domains_info")
	cursor, err := coll.Aggregate(context.TODO(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing aggregation: " + fmt.Sprint(err)})
		return
	}

	var result []bson.M
	if err = cursor.All(context.TODO(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding aggregation result: " + fmt.Sprint(err)})
		return
	}

	c.JSON(http.StatusOK, result)
}

/*
DeleteDomainById handles the DELETE request to remove a domain information register by its ID from MongoDB
Args:
	c *gin.Context: The Gin context for handling the request and response
Returns:
	None: Sends a JSON response confirming deletion or an error message
*/

func (h *Handler) DeleteDomainById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format: " + fmt.Sprint(err)})
		return
	}
	coll := h.DB.Client.Database(h.DB.DbName).Collection("domains_info")
	result, err := coll.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting domain: " + fmt.Sprint(err)})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Domain not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Domain info successfully deleted"})

}
