package response

import (
	"encoding/json"
	"net/http"
)

// Response encapsulate all response
type Response struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Info string      `json:"info"`
}

// JSON writes json http response
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
