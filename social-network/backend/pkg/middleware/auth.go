package middleware

import (
	"context"
	"net/http"

	"social-network/backend/pkg/repositories"
)

type ctxKey string

const userIDKey ctxKey = "userID"

// SessionMiddleware attaches user id to context when session is valid.
// It is permissive (does not block unauthenticated requests).
func SessionMiddleware(next http.Handler, sessionRepo repositories.SessionRepository, cookieName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err == nil {
			if uid, err := sessionRepo.LookupUserID(cookie.Value); err == nil && uid != "" {
				// store typed key to avoid collisions
				ctx := context.WithValue(r.Context(), userIDKey, uid)
				// keep string key for backward compatibility
				ctx = context.WithValue(ctx, "userID", uid)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAuth enforces that a request has an authenticated user in context.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if uid, ok := ctx.Value(userIDKey).(string); ok && uid != "" {
			next.ServeHTTP(w, r)
			return
		}
		if uid, ok := ctx.Value("userID").(string); ok && uid != "" {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})
}
