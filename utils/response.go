package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  bool        `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func SendResponse(w http.ResponseWriter, statusCode int, status bool, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(Response{
		Status:  status,
		Data:    data,
		Message: message,
	})
}

