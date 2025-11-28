package middleware

import (
	"net/http"
	"os"
	"pos-api/internal/lib"

	"github.com/go-chi/jwtauth/v5"
)

var TokenAuth *jwtauth.JWTAuth;

func InitJWT() {
	secretKey := os.Getenv("SECRET_KEY")
	TokenAuth = jwtauth.New("HS256", []byte(secretKey), nil)
}

type Payload struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Role string `json:"role"`
}

func CreateToken(p Payload) (string, error) {
	claims := map[string]any{
    "id":       p.Id,
    "username": p.Username,
    "role":     p.Role,
  }
	_, tokenString, err := TokenAuth.Encode(claims);

	if err != nil {
		return "", err;
	}

	return tokenString, nil;
}

func ExtractPayload(r *http.Request) (*Payload, error) {
	token, claims, err := jwtauth.FromContext(r.Context());
	if err != nil || token == nil || claims == nil {
    return nil, &lib.AppError{
      Message: "Unauthorized",
     	StatusCode: http.StatusUnauthorized,
    };
  }

  id, _ := claims["id"].(string);
  u, _ := claims["username"].(string);
  role, _ := claims["role"].(string);

  return &Payload{
    Id: id,
    Username: u,
    Role: role,
  }, nil
}
