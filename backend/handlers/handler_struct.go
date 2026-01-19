package handlers

import (
	"github.com/Nebula-Challenge/config"
	"github.com/Nebula-Challenge/scripts"
	"sync"
)

// Struct to hold scan request status and result
type ScanRequest struct {
	Status         string                     `json:"status"`
	Result         []byte                     `json:"result"` // Store the TLS assessment result
	FilteredResult *scripts.FilteredTLSReport `json:"filteredResult"`
	Error          string                     `json:"error"`
}

/*
Handler struct to hold the MongoDB client instance
*/
type Handler struct {
	DB           *config.DatabaseConfig  // The MongoDB Instance
	mu           sync.Mutex              // Mutex to protect access to scanRequests map
	scanRequests map[string]*ScanRequest // Map to store scan requests and their statuses

}

/*
NewHandler is used to create an instance of the handler Struct

params

	db *config.DatabaseConfig: pointer to the MongoDB configuration struct

return

	*Handler: pointer to a new Hanlder instance
*/
func NewHandler(db *config.DatabaseConfig) *Handler {
	return &Handler{
		DB:           db,
		scanRequests: make(map[string]*ScanRequest),
	}
}
