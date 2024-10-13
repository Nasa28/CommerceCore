package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DbUser                      string
	DBPassword                  string
	DBAddress                   string
	DBName                      string
	Port                        string
	PublicHost                  string
	JWTTokenExpirationInSeconds int64
	JWTSecret                      string
}

var Env = initConfig()

func initConfig() Config {
	// Load environment variables from .env file
	godotenv.Load()

	return Config{
		PublicHost:                  getEnv("PUBLIC_HOST", "http://localhost"),
		DbUser:                      getEnv("DB_USER", "root"),
		DBPassword:                  getEnv("DB_PASSWORD", "root"),
		DBAddress:                   fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		Port:                        getEnv("PORT", "5050"),
		DBName:                      getEnv("DB_NAME", "commercecome"),
		JWTTokenExpirationInSeconds: getEnvAsInt("JWTTokenExpirationInSeconds", 3600*24),
		JWTSecret:                      getEnv("JWT_SECRET", "secret"),
	}
}

// getEnv retrieves the value of the environment variable or uses a fallback if it's not set.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvAsInt retrieves the value of the environment variable as an int64 or returns the fallback.
func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return intValue
		}
	}
	return fallback
}
