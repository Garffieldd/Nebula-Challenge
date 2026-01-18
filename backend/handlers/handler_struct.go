package handlers

import (
	"sync"
	"github.com/Nebula-Challenge/config"

)

// Struct to hold scan request status and result
type ScanRequest struct {
    Status string `json:"status"`
    Result []byte `json:"result"`// Store the TLS assessment result
    Error  string `json:"error"` 
}

/*
Handler struct to hold the MongoDB client instance
*/
type Handler struct {
    DB           *config.DatabaseConfig      // The MongoDB Instance
    mu           sync.Mutex					// Mutex to protect access to scanRequests map		
    scanRequests map[string]*ScanRequest  // Map to store scan requests and their statuses

}

func NewHandler(db *config.DatabaseConfig) *Handler {
    return &Handler{
        DB:           db,
        scanRequests: make(map[string]*ScanRequest),
    }
}