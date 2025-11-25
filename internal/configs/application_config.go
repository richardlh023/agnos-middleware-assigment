package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ApplicationConfig struct {
	App struct {
		Port string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	JWT struct {
		Secret string
	}
	HISAPI struct {
		BaseURL string
	}
}

func LoadConfig() *ApplicationConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}

	config := &ApplicationConfig{}

	// Application Configuration
	config.App.Port = getEnv("SERVER_PORT", "8080")

	// Database Configuration
	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnv("DB_PORT", "5432")
	config.Database.User = getEnv("DB_USER", getEnv("USER", "postgres"))
	config.Database.Password = getEnv("DB_PASSWORD", "")
	config.Database.DBName = getEnv("DB_NAME", "agnos_db")

	// JWT Configuration
	config.JWT.Secret = getEnv("JWT_SECRET", "your-secret-key")

	// External HIS API Configuration
	config.HISAPI.BaseURL = getEnv("HIS_API_BASE_URL", "https://hospital-a.api.co.th")

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
