package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariables struct {
	Port        string
	DatabaseUrl string
}

func LoadEnvVariables() *EnvVariables {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DATABASE_URL")

	return &EnvVariables{
		Port:        port,
		DatabaseUrl: dbUrl,
	}
}
