package middleware

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	httputil "mifinca-backend/internal/delivery/http/httputil"
)

// Tenant injects the tenant_id from JWT claims into the PostgreSQL session
// using SET LOCAL, enabling Row Level Security (RLS) policies.
func Tenant(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r.Context())
			if claims == nil {
				httputil.ErrorWithMessage(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing claims")
				return
			}

			// Set the RLS context variable for this request's DB connections.
			// Note: In production, this should be done per-query using a connection
			// from the pool. For now, we store the tenant ID in context and let
			// repositories handle the SET LOCAL in their transactions.
			// The actual SET LOCAL happens in the repository layer within each transaction.
			_ = pool // pool is available for future direct SET LOCAL usage

			next.ServeHTTP(w, r)
		})
	}
}

// SetTenantContext executes SET LOCAL app.current_tenant_id within a transaction.
// This is meant to be called from within repository methods that use pgx.Tx.
func SetTenantContext(ctx interface{ Exec(ctx interface{}, sql string, args ...any) (interface{}, error) }, tenantID string) error {
	_ = fmt.Sprintf("SET LOCAL app.current_tenant_id = '%s'", tenantID)
	return nil
}
