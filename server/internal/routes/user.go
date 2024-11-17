package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/martishin/react-golang-oauth/internal/config"
	"github.com/martishin/react-golang-oauth/internal/handlers"
)

func RegisterUserRoutes(r *chi.Mux, app *config.Application) {
	userHandler := handlers.NewUserHandler(app.UserRepo)

	r.Group(func(protected chi.Router) {
		protected.Use(AuthMiddleware)
		protected.Get("/api/user", userHandler.User)
	})
}
