package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/martishin/react-golang-oauth/internal/config"
	"github.com/martishin/react-golang-oauth/internal/models"
)

func GenerateJWT(userID string, ttl time.Duration) (string, error) {
	claims := &models.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JwtKey)
}
