package main

import (
	"net/http"
)

type authenticationMiddleware struct {
	token string
}

func generateAuthenticationMiddleware() authenticationMiddleware {
	amw := authenticationMiddleware{}
	amw.generateToken()
	return amw
}

func (amw *authenticationMiddleware) generateToken() {
	amw.token = TOKEN
}

func (amw *authenticationMiddleware) authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")

		if token == amw.token {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
