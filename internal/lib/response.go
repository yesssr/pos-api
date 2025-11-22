package lib

import (
	"encoding/json"
	"net/http"
)

type Pagination struct {
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
}

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
	success bool,
	statusCode int,
	message string,
	data any,
	pagination *Pagination,
	credentials *string,
) {
	response := Response{
		Success:     success,
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

type AppError struct {
	Message string `json:"message"`
	StatusCode int `json:"status_code"`
}

func (e *AppError) Error() string {
	return e.Message;
}

func SendErrorResponse(
	w http.ResponseWriter,
	err error,
) {
	statusCode := http.StatusInternalServerError;

	if err, ok := err.(*AppError); ok {
		statusCode = err.StatusCode;
	}

	SendResponse(w, false, statusCode, err.Error(), nil, nil, nil)
}
