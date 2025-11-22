package lib

import (
	"os"

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
