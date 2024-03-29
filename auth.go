package auth

import (
	"context"
	"net/http"

	"github.com/paceperspective/auth/bajwt"
)

var ProjectID string

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString, err := bajwt.GetTokenFromHttpHeader(tokenHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		err = bajwt.Verify(r.Context(), ProjectID, tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CreateJWT creates a JWT with an additional username, uses a default expiry time of 1 hour
func CreateJWT(ctx context.Context, userName string) (string, error) {
	return bajwt.Create(ctx, ProjectID, userName, bajwt.StandardTokenLife)
}
