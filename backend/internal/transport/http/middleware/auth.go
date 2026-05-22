package middleware

import (
	"errors"
	"net/http"
	"strings"

	"gitflic.ru/lms/backend/internal/transport/http/response"
	"gitflic.ru/lms/backend/internal/transport/http/session"
)

// Auth verifies bearer session tokens and stores claims in request context.
func Auth(manager *session.Manager) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := bearerToken(r.Header.Get("Authorization"))
			if token == "" {
				response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "bearer token is required"))
				return
			}
			claims, err := manager.Verify(token)
			if err != nil {
				status := http.StatusUnauthorized
				code := "invalid_token"
				if errors.Is(err, session.ErrExpiredToken) {
					code = "expired_token"
				}
				response.WriteError(w, r, response.NewError(status, code, "session token is invalid"))
				return
			}
			next(w, r.WithContext(withClaims(r.Context(), claims)))
		}
	}
}

func bearerToken(header string) string {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return parts[1]
}
