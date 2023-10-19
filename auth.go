package auth

import (
	"net/http"

	"github.com/paceperspective/auth/bajwt"
)

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

		err = bajwt.Verify(r.Context(), tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
