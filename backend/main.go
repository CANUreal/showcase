package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"showcase/ent"
	storage "showcase/garage"
	"showcase/middlewares"
	"showcase/routes"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func openEntClient() (*ent.Client, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv)), nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("can't load dotenv %v", err)
	} 

	ctx := context.Background()

	entClient, err := openEntClient()
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer entClient.Close()

	garageClient, err := storage.NewGarageClient(ctx)
	if err != nil {
		log.Fatalf("failed creating garage client: %v", err)
	}

	profileImageHandler := &routes.ProfileImageHandler{Garage: garageClient, Ent: entClient}
	userHandler := &routes.UserHandler{Ent: entClient}
	authHandler := &routes.AuthHandler{Ent: entClient}
	requireSession := middlewares.RequireSession(entClient)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /projects", routes.Projects)
	mux.Handle("POST /users/profileupload", requireSession(profileImageHandler))

	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)

	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("GET /users", userHandler.List)
	mux.HandleFunc("GET /users/{id}", userHandler.Get)
	mux.Handle("PATCH /users/{id}", requireSession(http.HandlerFunc(userHandler.Update)))
	mux.Handle("DELETE /users/{id}", requireSession(http.HandlerFunc(userHandler.Delete)))

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
