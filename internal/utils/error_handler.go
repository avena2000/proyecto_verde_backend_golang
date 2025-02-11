package utils

import (
	"log"
	"net/http"
	"encoding/json"
)

func HandleError(w http.ResponseWriter, err error, status int, message string) {
	log.Printf("Error: %v - %s", err, message)
	
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
} 