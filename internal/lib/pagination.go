package lib

import (
	"context"
	"net/http"
	"strconv"
)

type Pagination struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalPages  *int `json:"total_pages,omitempty"`
}

type paginationKey struct{}

func Paginate(next http.Handler) http.Handler {
	hfn := func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page");
		perPageStr := r.URL.Query().Get("per_page");

		page := 1;
		perPage := 10;

		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p;
			}
		}

		if perPageStr != "" {
			if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 {
				perPage = pp;
			}
		}

		p := &Pagination{
			CurrentPage: page,
			PerPage:     perPage,
			TotalPages:  nil,
		}

		ctx := context.WithValue(r.Context(), paginationKey{}, p);
		next.ServeHTTP(w, r.WithContext(ctx));
	}
	return http.HandlerFunc(hfn);
}

func GetPagination(r *http.Request) *Pagination {
	p, _ := r.Context().Value(paginationKey{}).(*Pagination);
	return p;
}
