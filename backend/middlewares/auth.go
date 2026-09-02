package middlewares

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"showcase/ent"
	"showcase/ent/session"
	"time"

	"github.com/google/uuid"
)

const SessionCookieName = "session_token"

type ctxKey int

const userIDKey ctxKey = 0

func HashSessionToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}

// RequireSession validates the session cookie, checks it hasn't expired, and
// injects the owning user's ID into the request context for downstream handlers.
func RequireSession(entClient *ent.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil || cookie.Value == "" {
				http.Error(w, "not authenticated", http.StatusUnauthorized)
				return
			}

			tokenHash := HashSessionToken(cookie.Value)
			sess, err := entClient.Session.Query().
				Where(session.TokenHash(tokenHash)).
				WithOwner().
				Only(r.Context())
			if err != nil {
				if !ent.IsNotFound(err) {
					log.Printf("session lookup error: %v", err)
				}
				http.Error(w, "not authenticated", http.StatusUnauthorized)
				return
			}

			if sess.ExpiresAt.Before(time.Now()) {
				http.Error(w, "session expired", http.StatusUnauthorized)
				return
			}

			owner, err := sess.Edges.OwnerOrErr()
			if err != nil {
				log.Printf("session owner error: %v", err)
				http.Error(w, "not authenticated", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, owner.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
