package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"showcase/ent"
	"showcase/ent/schema"
	storage "showcase/garage"
	"showcase/middlewares"
)

const maxUploadSize = 5 << 20 // 5 megabytes

type ProfileImageHandler struct {
	Garage *storage.GarageClient
	Ent    *ent.Client
}

func (h *ProfileImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, ok := middlewares.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "file too large or invalid form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "missing 'image' field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/webp" && contentType != "image/png" {
		log.Printf("garage upload error %v", err)
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	key, err := h.Garage.UploadProfileImage(r.Context(), userId, file, contentType)
	if err != nil {
		log.Printf("garage upload error: %v", err)
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	_, err = h.Ent.User.UpdateOneID(userId).SetProfileLink(key).Save(r.Context())
	if err != nil {
		log.Printf("db update error: %v", err)
		http.Error(w, "failed to save profile link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"key": "key",
		"url": schema.ProfileImageURL(key),
	})
}
