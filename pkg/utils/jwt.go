package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

type JWTClaims struct {
	UserID  int    `json:"user_id"`
	Type    string `json:"type"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// generate JWT access
func GenerateAccessToken(UserID int, IsAdmin bool) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour) // berlaku 1 jam
	claims := &JWTClaims{
		UserID:  UserID,
		Type:    "access",
		IsAdmin: IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// generate JWT refresh
func GenerateRefreshTokenJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &JWTClaims{
		UserID: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// verify JWT
func VerifyJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
