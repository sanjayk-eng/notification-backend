package service

import (
	"fmt"
	"math/rand"
	"sanjay/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func GenerateAccessToken(phone string) (string, error) {
	claims := jwt.MapClaims{
		"phone": phone,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.LoadEnv().GetSecrateKey()))
}

// Generate Refresh Token (long-lived)
func GenerateRefreshToken(phone string) (string, error) {
	claims := jwt.MapClaims{
		"phone": phone,
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.LoadEnv().GetSecrateKey()))
}

// Validate any JWT
func ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.LoadEnv().GetSecrateKey()), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
