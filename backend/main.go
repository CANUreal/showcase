package main

import (
	"log"
	"net/http"
	"showcase/middlewares"
	"showcase/routes"
	"time"
)

type Resp struct {
	Status int `json:"status"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /projects", routes.Projects)

	handler := middlewares.LogMiddleware(mux)

	server := http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("Server started on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Printf("%v", err)
		return
	}
}
