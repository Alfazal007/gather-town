package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariables struct {
	Port               string
	DatabaseUrl        string
	AccessTokenSecret  string
	RefreshTokenSecret string
}

func LoadEnvVariables() *EnvVariables {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DATABASE_URL")
	accessToken := os.Getenv("ACCESS_TOKEN_SECRET")
	refreshToken := os.Getenv("REFRESH_TOKEN_SECRET")

	return &EnvVariables{
		Port:               port,
		DatabaseUrl:        dbUrl,
		AccessTokenSecret:  accessToken,
		RefreshTokenSecret: refreshToken,
	}
}
