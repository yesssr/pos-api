package middleware

import (
	"context"
	"net/http"
	"pos-api/internal/lib"
	"strings"
)

type queryKey struct{};
type QueryPayload struct {
	Search string
	OrderBy string
	OrderDir string
}

func QueryCtx(next http.Handler) http.Handler {
	hfn := func(w http.ResponseWriter, r *http.Request) {
	search := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("search")));
	orderBy := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("order_by")));
	orderDir := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("order_dir")));

		if orderBy != "" {
			allowedOrderBy := map[string]bool{
				"created_at": true,
				"username":   true,
				"role":       true,
				"price":			true,
				"stock":			true,
				"updated_at": true,
				"name":       true,
			}

			if !allowedOrderBy[orderBy] {
				lib.SendErrorResponse(w, &lib.AppError{
					Message: "Invalid order_by parameter",
					StatusCode: http.StatusBadRequest,
				}, nil);
				return;
			}
		} else {
			orderBy = "created_at";
		}

		if orderDir != "" {
			if orderDir != "asc" && orderDir != "desc" {
				lib.SendErrorResponse(w, &lib.AppError{
					Message: "Invalid order_dir parameter",
					StatusCode: http.StatusBadRequest,
				}, nil);
				return;
			}
		} else {
			orderDir = "desc";
		}

		if search != "" && len(search) < 3 {
			lib.SendErrorResponse(w, &lib.AppError{
				Message: "Search term must be at least 3 characters long",
				StatusCode: http.StatusBadRequest,
			}, nil);
			return;
		}

		q := &QueryPayload{
			Search: search,
			OrderBy: orderBy,
			OrderDir: orderDir,
		}

		ctx := context.WithValue(r.Context(), queryKey{}, q);
		next.ServeHTTP(w, r.WithContext(ctx));
	};
	return http.HandlerFunc(hfn);
}

func GetQueryFromCtx(r *http.Request) (*QueryPayload, bool) {
	s, ok := r.Context().Value(queryKey{}).(*QueryPayload);
	return s, ok;
}
