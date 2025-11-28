package middleware

import (
	"context"
	"net/http"
	"pos-api/internal/lib"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

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
