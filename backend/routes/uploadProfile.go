package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"showcase/ent"
	"showcase/ent/schema"
	storage "showcase/garage"

	"github.com/google/uuid"
) 

const maxUploadSize = 5 << 20 // 5 megabytes

type ProfileImageHandler struct {
	Garage	*storage.GarageClient
	Ent 	*ent.Client
}

func (h *ProfileImageHandler) uploadProfileRoute(w http.ResponseWriter, r *http.Request) {
	// NOTICE !
	// WHEN WE COMPLETE AUTH MIDDLEWARE
	// WE WILL USE THE USERID FROM THERE
	// DO NOT FORGET THIS
	// WHEN WE COMPLETE IT DELETE THESE COMMENTS
	// FOR NOW WE WILL TEST IT WITH QUERY PARAMS
	
	userIdStr := r.URL.Query().Get("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
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
