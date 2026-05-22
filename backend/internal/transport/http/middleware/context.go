package middleware

import (
	"context"

	"gitflic.ru/lms/backend/internal/transport/http/session"
)

type claimsContextKey struct{}
type requestIDContextKey struct{}

// ClaimsFromContext returns authenticated session claims.
func ClaimsFromContext(ctx context.Context) (session.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(session.Claims)
	return claims, ok
}

func withClaims(ctx context.Context, claims session.Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey{}, claims)
}

// RequestIDFromContext returns the request ID assigned by middleware.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDContextKey{}).(string)
	return requestID, ok
}

func withRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey{}, requestID)
}
