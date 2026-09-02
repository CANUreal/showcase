package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"showcase/ent"
	"showcase/middlewares"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Ent *ent.Client
}

type createUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type userResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	ProfileLink string    `json:"profile_link"`
}

func toUserResponse(u *ent.User) userResponse {
	return userResponse{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		ProfileLink: u.ProfileLink,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "username, email and password are required", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("password hash error: %v", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	user, err := h.Ent.User.Create().
		SetUsername(req.Username).
		SetEmail(req.Email).
		SetPasswordHash(string(hash)).
		Save(r.Context())
	if err != nil {
		if ent.IsConstraintError(err) {
			http.Error(w, "username or email already taken", http.StatusConflict)
			return
		}
		if ent.IsValidationError(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("db create error: %v", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toUserResponse(user))
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	user, err := h.Ent.User.Get(r.Context(), userId)
	if err != nil {
		if ent.IsNotFound(err) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		log.Printf("db get error: %v", err)
		http.Error(w, "failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toUserResponse(user))
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.Ent.User.Query().All(r.Context())
	if err != nil {
		log.Printf("db list error: %v", err)
		http.Error(w, "failed to list users", http.StatusInternalServerError)
		return
	}

	resp := make([]userResponse, len(users))
	for i, u := range users {
		resp[i] = toUserResponse(u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	sessionUserId, ok := middlewares.UserIDFromContext(r.Context())
	if !ok || sessionUserId != userId {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	update := h.Ent.User.UpdateOneID(userId).
		SetNillableUsername(req.Username).
		SetNillableEmail(req.Email)

	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("password hash error: %v", err)
			http.Error(w, "failed to update user", http.StatusInternalServerError)
			return
		}
		update = update.SetPasswordHash(string(hash))
	}

	user, err := update.Save(r.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		if ent.IsConstraintError(err) {
			http.Error(w, "username or email already taken", http.StatusConflict)
			return
		}
		if ent.IsValidationError(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("db update error: %v", err)
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toUserResponse(user))
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	sessionUserId, ok := middlewares.UserIDFromContext(r.Context())
	if !ok || sessionUserId != userId {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := h.Ent.User.DeleteOneID(userId).Exec(r.Context()); err != nil {
		if ent.IsNotFound(err) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		log.Printf("db delete error: %v", err)
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
