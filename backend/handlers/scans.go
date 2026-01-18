package handlers

import (
	"net/http"

	"github.com/Nebula-Challenge/scripts"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
StartScan handles the POST request to start a TLS scan for a given domain
Args:

	c *gin.Context: The Gin context for handling the request and response

Returns:

	None: Sends a JSON response with the scan request ID
*/
func (h *Handler) StartScan(c *gin.Context) {
	var req struct {
		Domain string `json:"domain"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	scanRequestID := uuid.New().String()
	h.mu.Lock() //Acceder a la gorutina de manera segura
	h.scanRequests[scanRequestID] = &ScanRequest{Status: "IN_PROGRESS"}
	h.mu.Unlock()

	go func() { // gorutina para manejar la evaluacion asincronamente (un hilo ligero de go)
		_, err := scripts.CheckTLS(req.Domain, true)
		if err != nil {
			h.updateScanRequest(scanRequestID, "error", nil, err.Error())
			return
		}

		result, err := scripts.PollUntilReady(req.Domain)
		if err != nil {
			h.updateScanRequest(scanRequestID, "error", nil, err.Error())
		} else {
			h.updateScanRequest(scanRequestID, "complete", result, "")
		}

	}()
	c.JSON(http.StatusOK, gin.H{"scanRequestID": scanRequestID})
}

/*
GetScanStatus handles the GET request to retrieve the status of a TLS scan, its porpous is to give the frontend a way to check the status of a previously initiated scan and ,
if completed, retrieve the results.
Args:

	c *gin.Context: The Gin context for handling the request and response

Returns:

	None: Sends a JSON response with the scan status and result or an error message
*/
func (h *Handler) GetScanStatus(c *gin.Context) {
	scanRequestID := c.Param("scanRequestID")
	h.mu.Lock()
	scanRequest, exists := h.scanRequests[scanRequestID]
	h.mu.Unlock()
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scan request not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": scanRequest.Status, "result": scanRequest.Result, "error": scanRequest.Error})
}

/*
updateScanRequest updates the status and result of a scan request map in a thread-safe manner
Args:

	id string: The scan request ID
	status string: The new status of the scan request
	result map[string]interface{}: The result of the scan request
	errMsg string: Any error message associated with the scan request
*/
func (h *Handler) updateScanRequest(id string, status string, result []byte, errMsg string) {
	h.mu.Lock()
	defer h.mu.Unlock()                            // Para evitar condidiones de carrera
	if value, exist := h.scanRequests[id]; exist { // Verifica que el ID exista antes de actualizar el estado de este scan Request
		value.Status = status
		value.Result = result
		value.Error = errMsg
	}
}
