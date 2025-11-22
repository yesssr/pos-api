package lib

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
  validate = validator.New()
}

func ValidateStruct(s any) error {
  return validate.Struct(s)
}

func ValidateJSON[T any](w http.ResponseWriter, r *http.Request, dst *T) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		SendErrorResponse(w, &AppError{
      Message: "Invalid JSON",
      StatusCode: http.StatusBadRequest,
    })
		return false
	}
	return true
}
