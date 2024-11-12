package utils

import (
	"time"

	"github.com/Alfazal007/gather-town/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(user database.User) (string, string, error) {
	accessToken, err := generateJWT(user)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func generateJWT(user database.User) (string, error) {
	jwtSecret := LoadEnvVariables().AccessTokenSecret

	secretKey := []byte(jwtSecret)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["username"] = user.Username

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func generateRefreshToken(user database.User) (string, error) {
	jwtSecret := LoadEnvVariables().RefreshTokenSecret
	secretKey := []byte(jwtSecret)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["username"] = user.Username
	claims["email"] = user.Email

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
