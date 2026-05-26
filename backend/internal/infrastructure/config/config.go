package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration parameters for the application.
type Config struct {
	Port                string
	Env                 string
	DatabaseURL         string
	GeminiAPIKey        string
	AdminAPIKey         string
	CronFetchSchedule   string
	CronProcessSchedule string
	FrontendURL         string
	AIAPIKey            string
	AIBaseURL           string
	AIModel             string
}

// Load reads configuration variables from environment variables with sensible defaults.
func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error reading it, relying on existing environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	cronFetch := os.Getenv("CRON_FETCH_SCHEDULE")
	if cronFetch == "" {
		cronFetch = "0 9 * * *" // Daily at 9:00 AM
	}

	cronProcess := os.Getenv("CRON_PROCESS_SCHEDULE")
	if cronProcess == "" {
		cronProcess = "0 10 * * *" // Daily at 10:00 AM
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	aiBaseURL := os.Getenv("AI_BASE_URL")
	if aiBaseURL == "" {
		aiBaseURL = "https://api.mistral.ai/v1"
	}

	aiModel := os.Getenv("AI_MODEL")
	if aiModel == "" {
		aiModel = "mistral-large-latest"
	}

	return Config{
		Port:                port,
		Env:                 env,
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		GeminiAPIKey:        os.Getenv("GEMINI_API_KEY"),
		AdminAPIKey:         os.Getenv("ADMIN_API_KEY"),
		CronFetchSchedule:   cronFetch,
		CronProcessSchedule: cronProcess,
		FrontendURL:         frontendURL,
		AIAPIKey:            os.Getenv("AI_API_KEY"),
		AIBaseURL:           aiBaseURL,
		AIModel:             aiModel,
	}
}
