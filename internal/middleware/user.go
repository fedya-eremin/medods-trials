package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/fedya-eremin/medods-trials/internal/contextkeys"
	"github.com/fedya-eremin/medods-trials/internal/service/jwt"
)

func JWTMiddleware(j *jwt.JWTService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Invalid auth header", http.StatusUnauthorized)
				return
			}
			claims, err := j.ParseAccessToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}
			if time.Now().After(claims.ExpiresAt.Time) {
				http.Error(w, "Access token expired", http.StatusUnauthorized)
			}
			ctx := contextkeys.WithContextValue(r.Context(), contextkeys.JWTClaimsKey, claims)
			next(w, r.WithContext(ctx))
		})
	}
}
