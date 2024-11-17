package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/martishin/react-golang-oauth/internal/config"
	"github.com/martishin/react-golang-oauth/internal/handlers"
)

func RegisterAuthRoutes(r *chi.Mux, app *config.Application) {
	authHandler := handlers.NewAuthHandler(app.UserRepo)

	r.Post("/api/auth/google", authHandler.GoogleLogin)
	r.Post("/api/auth/logout", authHandler.Logout)
	r.Post("/api/auth/refresh", authHandler.RefreshToken)
}
