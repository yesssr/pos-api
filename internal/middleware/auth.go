package middleware

import (
	"context"
	"net/http"
	"pos-api/internal/lib"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type payloadKey struct {};
var userPayload payloadKey;

func Auth(next http.Handler) http.Handler {
  hfn := func(w http.ResponseWriter, r *http.Request) {
  	p, err := lib.ExtractPayload(r);

   	if err != nil {
    	lib.SendErrorResponse(w, err, nil);
     	return;
    }

  	ctx := context.WithValue(r.Context(), userPayload, p);
  	next.ServeHTTP(w, r.WithContext(ctx));
  }
  return http.HandlerFunc(hfn);
}


func GetUserPayload(r *http.Request) (*lib.Payload, error) {
	p, ok := r.Context().Value(userPayload).(*lib.Payload);
	if !ok || p == nil {
		return nil, &lib.AppError{
			Message: "Unauthorized",
			StatusCode: http.StatusUnauthorized,
		};
	}
	return p, nil;
}


func IsAdmin(next http.Handler) http.Handler {
  hfn := func(w http.ResponseWriter, r *http.Request) {
    p, err := GetUserPayload(r);
    if err != nil {
			lib.SendErrorResponse(w, err, nil)
    }

    if p.Role != "admin" {
    	lib.SendErrorResponse(w, &lib.AppError{
		 		Message: "Forbidden",
			 	StatusCode: http.StatusForbidden,
     	}, nil);
			return;
    }
    next.ServeHTTP(w, r);
  }
  return http.HandlerFunc(hfn);
}

type idKey struct{};

func IdCtx(next http.Handler) http.Handler {
	hfn := func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id");

		id, err := uuid.Parse(idStr);
	 	if err != nil {
			lib.SendErrorResponse(w, &lib.AppError{
				Message: "Invalid ID format",
				StatusCode: http.StatusBadRequest,
			}, nil);
			return;
		}

		pgId := pgtype.UUID{
			Bytes: id,
			Valid: true,
		}

		ctx := context.WithValue(r.Context(), idKey{}, pgId);
		next.ServeHTTP(w, r.WithContext(ctx));
	}

	return http.HandlerFunc(hfn);
}

func GetIdFromCtx(r *http.Request) (pgtype.UUID, error) {
	id, ok := r.Context().Value(idKey{}).(pgtype.UUID);
	if !ok {
		return pgtype.UUID{}, &lib.AppError{
			Message: "ID not found in context",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return id, nil;
}
