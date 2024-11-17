package config

import "github.com/martishin/react-golang-oauth/internal/repository"

var JwtKey = []byte("YOUR_GOOGLE_APP_CLIENT_SECRET")

type Application struct {
	UserRepo repository.UserRepository
}
