package lib

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success     bool        `json:"success"`
	StatusCode  int         `json:"status_code"`
	Message     string      `json:"message"`
	Data        any         `json:"data,omitempty"`
	Pagination  *Pagination `json:"pagination,omitempty"`
	Credentials *string     `json:"credentials,omitempty"`
}

func SendResponse(
	w http.ResponseWriter,
	statusCode int,
	message string,
	data any,
	pagination *Pagination,
	credentials *string,
) {
	response := Response{
		Success:     true,
		StatusCode:  statusCode,
		Message:     message,
		Data:        data,
		Pagination:  pagination,
		Credentials: credentials,
	}
	w.Header().Set("Content-Type", "application/json");
	w.WriteHeader(statusCode);
	json.NewEncoder(w).Encode(response);
}
