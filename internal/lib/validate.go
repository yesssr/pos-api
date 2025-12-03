package lib

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

var validate *validator.Validate
func usernameValidator(fl validator.FieldLevel) bool {
  u := fl.Field().String()
  if u != strings.ToLower(u) {
    return false;
  }
  if strings.Contains(u, " ") {
    return false;
  }
  return true;
}

func init() {
  validate = validator.New()
  validate.RegisterValidation("username", usernameValidator);
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
    }, nil)
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

func IntToPgNumeric(v int) *pgtype.Numeric {
	r := big.NewInt(int64(v));
	return &pgtype.Numeric{
		Int: r,
		Valid: true,
	}
}

func NumericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}

	i := new(big.Float).SetInt(n.Int)
	e := big.NewFloat(math.Pow10(int(n.Exp)))
	f := new(big.Float).Mul(i, e)

	v, _ := f.Float64()
	return v
}

func GenerateUniqueNumber() string {
	now := time.Now();
	datePart := now.Format("20060102");

	r := rand.New(rand.NewSource(time.Now().UnixNano()));
	randomPart := r.Intn(9000) + 1000;

	uniqueNumber := fmt.Sprintf("%s%d", datePart, randomPart);
	return uniqueNumber;
}
