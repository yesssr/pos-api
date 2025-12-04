package middleware

import (
	"context"
	"net/http"
	"pos-api/internal/lib"
	"strings"
	"time"
)

type queryKey struct{};
type QueryPayload struct {
	Search string;
	OrderBy string;
	OrderDir string;
	Period string;
	StartAt time.Time;
	EndAt time.Time;
}

func QueryCtx(allowedCol, allowedPer map[string]bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			search := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("search")));
			orderBy := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("order_by")));
			orderDir := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("order_dir")));
			period := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("period")));
			startAtStr := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("start_at")));
			endAtStr := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("end_at")));

			now := time.Now();
			startAt := now;
			endAt := now;

			if orderBy != "" {
				if !allowedCol[orderBy] {
					lib.SendErrorResponse(w, &lib.AppError{
						Message:    "Invalid order_by parameter",
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
						Message:    "Invalid order_dir parameter",
						StatusCode: http.StatusBadRequest,
					}, nil);
					return;
				}
			} else {
				orderDir = "desc";
			}

			if search != "" && len(search) < 3 {
				lib.SendErrorResponse(w, &lib.AppError{
					Message:    "Search term must be at least 3 characters long",
					StatusCode: http.StatusBadRequest,
				}, nil);
				return;
			}

			if startAtStr != "" && endAtStr != "" {
				layout := "2006-01-02";
				s, err := time.Parse(layout, startAtStr); if err != nil {
					lib.SendErrorResponse(w, &lib.AppError{
						Message:    "Invalid start_at parameter",
						StatusCode: http.StatusBadRequest,
					}, nil);
					return;
				}

				e, err := time.Parse(layout, endAtStr); if err != nil {
					lib.SendErrorResponse(w, &lib.AppError{
						Message:    "Invalid end_at parameter",
						StatusCode: http.StatusBadRequest,
					}, nil);
					return;
				}

				startAt = s;
				endAt = e;
			}

			if period != "" {
				if !allowedPer[period] {
					lib.SendErrorResponse(w, &lib.AppError{
						Message:    "Invalid period parameter",
						StatusCode: http.StatusBadRequest,
					}, nil);
					return;
				}
			} else {
				period = "day";
			}

			q := &QueryPayload{
				Search:   search,
				OrderBy:  orderBy,
				OrderDir: orderDir,
				StartAt:  startAt,
				EndAt:    endAt,
				Period:   period,
			}

			ctx := context.WithValue(r.Context(), queryKey{}, q)
			next.ServeHTTP(w, r.WithContext(ctx));
		});
	}
}


func GetQueryFromCtx(r *http.Request) (*QueryPayload, bool) {
	s, ok := r.Context().Value(queryKey{}).(*QueryPayload);
	return s, ok;
}
