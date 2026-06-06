package middleware

import (
	"context"
	"net/http"
	"strings"

	infraAuth "mifinca-backend/internal/infrastructure/auth"
	httputil "mifinca-backend/internal/delivery/http/httputil"
)

type contextKey string

const ClaimsKey contextKey = "claims"

// Auth returns a middleware that validates JWT tokens and injects claims into context.
func Auth(jwtMaker *infraAuth.JWTMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httputil.ErrorWithMessage(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				httputil.ErrorWithMessage(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid authorization format")
				return
			}

			claims, err := jwtMaker.VerifyToken(parts[1])
			if err != nil {
				httputil.ErrorWithMessage(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaims extracts the authenticated claims from the request context.
func GetClaims(ctx context.Context) *infraAuth.Claims {
	claims, ok := ctx.Value(ClaimsKey).(*infraAuth.Claims)
	if !ok {
		return nil
	}
	return claims
}
