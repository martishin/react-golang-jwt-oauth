package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/martishin/react-golang-oauth/internal/config"
	"github.com/martishin/react-golang-oauth/internal/models"
	"github.com/martishin/react-golang-oauth/internal/repository"
	"github.com/martishin/react-golang-oauth/internal/utils"
)

var authTokenDuration = 15 * time.Minute
var refreshTokenDuration = 7 * 24 * time.Hour

type AuthHandler struct {
	userRepo repository.UserRepository
}

func NewAuthHandler(repo repository.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: repo}
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IDToken string `json:"id_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		fmt.Printf("Error decoding request body: %v\n", err)
		return
	}

	userInfo, err := utils.ValidateGoogleIDToken(r.Context(), body.IDToken)
	if err != nil {
		http.Error(w, "Invalid Google ID Token", http.StatusUnauthorized)
		fmt.Printf("Error validating ID token: %v\n", err)
		return
	}

	userID, err := h.userRepo.FindOrCreateUser(r.Context(), userInfo)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Printf("Error accessing the database: %v\n", err)
		return
	}

	authToken, err := utils.GenerateJWT(userID, authTokenDuration)
	if err != nil {
		http.Error(w, "Failed to generate auth token", http.StatusInternalServerError)
		fmt.Printf("Error generating auth token: %v\n", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    authToken,
		Expires:  time.Now().Add(authTokenDuration),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	refreshToken, err := utils.GenerateJWT(userID, refreshTokenDuration)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		fmt.Printf("Error generating refresh token: %v\n", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(refreshTokenDuration),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	refreshToken := refreshTokenCookie.Value

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	newAuthToken, err := utils.GenerateJWT(claims.UserID, authTokenDuration)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := utils.GenerateJWT(claims.UserID, refreshTokenDuration)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    newAuthToken,
		Expires:  time.Now().Add(authTokenDuration),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Expires:  time.Now().Add(refreshTokenDuration),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
}
