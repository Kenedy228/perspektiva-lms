package middleware

import (
	"log/slog"
	"net/http"

	"gitflic.ru/lms/backend/internal/transport/http/response"
)

// Recover converts panics to a stable internal error response.
func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error("http panic", "panic", recovered, "path", r.URL.Path)
					response.WriteError(w, r, response.NewError(http.StatusInternalServerError, "internal_error", "internal server error"))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
