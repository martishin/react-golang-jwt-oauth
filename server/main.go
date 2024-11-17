package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/martishin/react-golang-oauth/internal/config"
	"github.com/martishin/react-golang-oauth/internal/db"
	"github.com/martishin/react-golang-oauth/internal/repository"
	"github.com/martishin/react-golang-oauth/internal/routes"
)

func main() {
	err := db.Connect("postgres://postgres:postgres@localhost:5432/oauth")
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v", err))
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

	userRepo := repository.NewUserRepository(db.GetDB())

	app := &config.Application{
		UserRepo: userRepo,
	}

	routes.RegisterAuthRoutes(r, app)
	routes.RegisterUserRoutes(r, app)

	fmt.Println("Server is running on :8000")
	http.ListenAndServe(":8000", r)
}
