package utils

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/oauth2/v2"
)

func ValidateGoogleIDToken(ctx context.Context, idToken string) (*oauth2.Userinfo, error) {
	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		certURL := "https://www.googleapis.com/oauth2/v3/certs"
		resp, err := http.Get(certURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Google certs: %v", err)
		}
		defer resp.Body.Close()

		var certs struct {
			Keys []struct {
				Kid string `json:"kid"`
				N   string `json:"n"`
				E   string `json:"e"`
			} `json:"keys"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&certs); err != nil {
			return nil, fmt.Errorf("failed to parse Google certs: %v", err)
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("invalid token header")
		}

		for _, key := range certs.Keys {
			if key.Kid == kid {
				modulus, err := base64.RawURLEncoding.DecodeString(key.N)
				if err != nil {
					return nil, fmt.Errorf("failed to decode modulus: %v", err)
				}
				exponent, err := base64.RawURLEncoding.DecodeString(key.E)
				if err != nil {
					return nil, fmt.Errorf("failed to decode exponent: %v", err)
				}

				e := 0
				for _, b := range exponent {
					e = e*256 + int(b)
				}

				pubKey := &rsa.PublicKey{
					N: new(big.Int).SetBytes(modulus),
					E: e,
				}
				return pubKey, nil
			}
		}

		return nil, errors.New("public key not found")
	})

	if err != nil {
		fmt.Printf("Error parsing ID token: %v\n", err)
		return nil, fmt.Errorf("invalid ID token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Printf("Error parsing token claims\n")
		return nil, errors.New("failed to parse token claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		fmt.Printf("Missing email claim\n")
		return nil, errors.New("missing email claim")
	}
	name, _ := claims["name"].(string)
	picture, _ := claims["picture"].(string)

	return &oauth2.Userinfo{
		Email:   email,
		Name:    name,
		Picture: picture,
	}, nil
}
