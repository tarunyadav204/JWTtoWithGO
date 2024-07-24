package utils

import (
	"jwtAuth/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("TARUN1508")

// GenerateToken generates a new JWT token for the given username
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &model.Claims{
		UserEmail: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyToken checks if the provided JWT token is valid and returns the claims
func VerifyToken(tokenStr string) (*model.Claims, error) {
	claims := &model.Claims{}

	// Parse the JWT token string `tokenStr` and store its claims into `claims`
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// In this function, you return the key used to sign the token
		return jwtKey, nil
	})

	// Check if there was an error parsing the token or if the token is not valid
	if err != nil || !token.Valid {
		return nil, err
	}

	// Return the parsed claims if the token is valid
	return claims, nil
}
