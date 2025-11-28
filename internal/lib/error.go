package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

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
		return fmt.Sprintf("Minimal %s karakter", fe.Param());
	case "max":
		return fmt.Sprintf("Maksimal %s karakter", fe.Param());
	case "oneof":
		return fmt.Sprintf("Harus salah satu dari: %s", fe.Param());
	case "username":
		return "Username harus huruf kecil tanpa spasi";
	}
	return "Format tidak valid"
}

func SendErrorResponse(
	w http.ResponseWriter,
	err error,
	obj any,
) {
	statusCode := http.StatusInternalServerError;
	msg := "Internal Server Error";
	var errors any = nil;
	fmt.Println(err);
	if err, ok := err.(*AppError); ok {
		statusCode = err.StatusCode;
		msg = err.Message;
	}

	if err, ok := err.(validator.ValidationErrors); ok {
		statusCode = http.StatusBadRequest;
		msg = "Validation Error";
		out := make(map[string]string);
  	typ := reflect.TypeOf(obj).Elem();

		for _, e := range err {
			fieldStruct, _ := typ.FieldByName(e.StructField());
			jsonTag := fieldStruct.Tag.Get("json");
			if jsonTag == "" {
				jsonTag = e.Field();
			}
			out[jsonTag] = validationMessage(e);
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
