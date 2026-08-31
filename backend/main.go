package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Resp struct {
	Status	int `json:"status"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(Resp{
			Status: 200,
		})
	})

	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	log.Printf("Server started on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Printf("%v", err)
		return
	}
}
