package lib

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

var validate *validator.Validate

func init() {
  validate = validator.New()
}

func ValidateStruct(s any) error {
  return validate.Struct(s)
}

func ValidateJSON[T any](w http.ResponseWriter, r *http.Request, dst *T) bool {
	fmt.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		fmt.Println(err);
		SendErrorResponse(w, &AppError{
      Message: "Invalid JSON",
      StatusCode: http.StatusBadRequest,
    })
		return false
	}
	return true
}

func BoolPtrToPgBool(b *bool,) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{
			Bool:  false,
			Valid: false,
		}
	}
	return pgtype.Bool{
		Bool:  *b,
		Valid: true,
	}
}

func IntToPgNumeric(v *big.Int) *pgtype.Numeric {
	return &pgtype.Numeric{
		Int: v,
		Valid: true,
	}
}

func GenerateUniqueNumber() string {
	now := time.Now();
	datePart := now.Format("20060102");

	r := rand.New(rand.NewSource(time.Now().UnixNano()));
	randomPart := r.Intn(9000) + 1000;

	uniqueNumber := fmt.Sprintf("%s%d", datePart, randomPart);
	return uniqueNumber;
}
