package routes

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"showcase/ent"
	"showcase/ent/session"
	"showcase/ent/user"
	"showcase/middlewares"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const sessionDuration = 7 * 24 * time.Hour

type AuthHandler struct {
	Ent *ent.Client
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func newSessionToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func setSessionCookie(w http.ResponseWriter, r *http.Request, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     middlewares.SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	foundUser, err := h.Ent.User.Query().Where(user.Username(req.Username)).Only(r.Context())
	if err != nil {
		if !ent.IsNotFound(err) {
			log.Printf("login lookup error: %v", err)
		}
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := newSessionToken()
	if err != nil {
		log.Printf("session token error: %v", err)
		http.Error(w, "failed to log in", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(sessionDuration)
	_, err = h.Ent.Session.Create().
		SetTokenHash(middlewares.HashSessionToken(token)).
		SetExpiresAt(expiresAt).
		SetOwnerID(foundUser.ID).
		Save(r.Context())
	if err != nil {
		log.Printf("session create error: %v", err)
		http.Error(w, "failed to log in", http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, r, token, expiresAt)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(middlewares.SessionCookieName)
	if err == nil && cookie.Value != "" {
		_, err := h.Ent.Session.Delete().
			Where(session.TokenHash(middlewares.HashSessionToken(cookie.Value))).
			Exec(r.Context())
		if err != nil {
			log.Printf("session delete error: %v", err)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middlewares.SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusNoContent)
}
