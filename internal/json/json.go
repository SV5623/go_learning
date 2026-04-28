package json

import (
	"encoding/json"
	"net/http"
)


func WriteJSON (w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	// 2 return response json in http
	json.NewEncoder(w).Encode(data)
} 