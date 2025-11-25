package service

import (
	"context"
	"net/http"
	"pos-api/internal/lib"
	"pos-api/internal/store"

	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	q *store.Queries;
}

func NewAuthService(q *store.Queries) *AuthService {
	return &AuthService{ q: q };
}

func(s *AuthService) Login(ctx context.Context, username, password string) (store.GetUserByUsernameRow, string, error) {
	u, err := s.q.GetUserByUsername(ctx, username);
	if err != nil {
		if err == pgx.ErrNoRows{
			return store.GetUserByUsernameRow{}, "", &lib.AppError{
				Message: "Invalid username or password",
				StatusCode: http.StatusUnauthorized,
			}
		}
		return store.GetUserByUsernameRow{}, "", err;
	}

	if err := lib.VerifyPassword(u.Password, password); err != nil {
		return store.GetUserByUsernameRow{}, "", &lib.AppError{
			Message: "Invalid username or password",
			StatusCode: http.StatusUnauthorized,
		};
	}

	p := lib.Payload{
		Id: u.ID.String(),
		Username: u.Username,
		Role: string(u.Role),
	}
	c, err := lib.CreateToken(p);
	if err != nil {
		return store.GetUserByUsernameRow{}, "", err;
	}
	return u, c, nil;
}
