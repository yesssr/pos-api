package lib

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AppError struct {
	Message string `json:"message"`
	StatusCode int `json:"status_code"`
}

func (e *AppError) Error() string {
	return e.Message;
}

type ErrorResponse struct {
	Success     bool        `json:"success"`
	StatusCode  int         `json:"status_code"`
	Message     string      `json:"message"`
	Errors        any       `json:"errors,omitempty"`
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Field ini wajib diisi"
	case "min":
		return fmt.Sprintf("Minimal %s karakter", fe.Param())
	case "max":
		return fmt.Sprintf("Maksimal %s karakter", fe.Param())
	case "oneof":
		return fmt.Sprintf("Harus salah satu dari: %s", fe.Param())
	case "email":
		return "Format email tidak valid"
	}
	return "Format tidak valid"
}

func SendErrorResponse(
	w http.ResponseWriter,
	err error,
) {
	statusCode := http.StatusInternalServerError;
	msg := "Internal Server Error";
	var errors any = nil;

	if err, ok := err.(*AppError); ok {
		statusCode = err.StatusCode;
		msg = err.Message;
	}

	if err, ok := err.(validator.ValidationErrors); ok {
		statusCode = http.StatusBadRequest;
		msg = "Validation Error";
		var out map[string]string = make(map[string]string);
		for _, e := range err {
			out[e.Field()] = validationMessage(e);
		}
		errors = out;
	}

	response := ErrorResponse{
		Success:    false,
		StatusCode: statusCode,
		Message:    msg,
		Errors:     errors,
	}

	w.Header().Set("Content-Type", "application/json");
	w.WriteHeader(statusCode);
	json.NewEncoder(w).Encode(response);
}
