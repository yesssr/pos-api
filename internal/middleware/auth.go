package middleware

import (
	"context"
	"net/http"
	"pos-api/internal/lib"
)

type payloadKey struct {};
var userPayload payloadKey;

func Auth() func(http.Handler) http.Handler {
  return func(next http.Handler) http.Handler {
  	hfn := func(w http.ResponseWriter, r *http.Request) {
   	p, err := lib.ExtractPayload(r);

    if err != nil {
      lib.SendErrorResponse(w, err);
      return;
    }

    ctx := context.WithValue(r.Context(), userPayload, p);
    next.ServeHTTP(w, r.WithContext(ctx));
  }
    return http.HandlerFunc(hfn);
  }
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

func IsAdmin() func(http.Handler) http.Handler {
  return func(next http.Handler) http.Handler {
  	hfn := func(w http.ResponseWriter, r *http.Request) {
    	p, err := GetUserPayload(r);
     	if err != nil {
				lib.SendErrorResponse(w, err)
      }

     	if p.Role != "admin" {
      	lib.SendErrorResponse(w, &lib.AppError{
			 		Message: "Forbidden",
				 	StatusCode: http.StatusForbidden,
       	})
				return;
      }
      next.ServeHTTP(w, r);
    }
    return http.HandlerFunc(hfn);
  }
}
