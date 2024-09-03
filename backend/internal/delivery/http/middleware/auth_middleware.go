package middleware

import (
	"context"
	"goP2Pbackend/internal/domain"
	"net/http"
	"strings"
)

func AuthMiddleware(userUsecase domain.UserUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 {
				http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := bearerToken[1]

			// TODO: Implement proper token validation
			// This is a placeholder implementation
			user, err := userUsecase.GetByID(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
