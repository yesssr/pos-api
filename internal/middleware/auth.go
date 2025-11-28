package middleware

import (
	"context"
	"net/http"
	"pos-api/internal/lib"
)

type payloadKey struct {};

func Auth(next http.Handler) http.Handler {
  hfn := func(w http.ResponseWriter, r *http.Request) {
  	p, err := ExtractPayload(r);

   	if err != nil {
    	lib.SendErrorResponse(w, err, nil);
     	return;
    }

  	ctx := context.WithValue(r.Context(), payloadKey{}, p);
  	next.ServeHTTP(w, r.WithContext(ctx));
  }
  return http.HandlerFunc(hfn);
}

func GetUserPayload(r *http.Request) (*Payload, error) {
	p, ok := r.Context().Value(payloadKey{}).(*Payload);
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
