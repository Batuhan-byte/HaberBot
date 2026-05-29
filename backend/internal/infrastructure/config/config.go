package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// mustEnv reads a required environment variable and fatally exits if not set.
func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("FATAL: Required environment variable %q is not set. Application cannot start.", key)
	}
	return val
}

// Config holds all configuration parameters for the application.
type Config struct {
	Port                 string
	Env                  string
	IsProduction         bool   // true when ENV=production; enables Secure cookies, strict headers etc.
	DatabaseURL          string
	GeminiAPIKey         string
	AdminAPIKey          string
	CronFetchSchedule    string
	CronProcessSchedule  string
	FrontendURL          string
	AIAPIKey             string
	AIBaseURL            string
	AIModel              string
	PendingLimitPerTopic int
	JWTSecret            string
	JWTAccessTTLMinutes  int
	JWTRefreshTTLDays    int
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

	pendingLimitStr := os.Getenv("PENDING_LIMIT_PER_TOPIC")
	pendingLimit := 50
	if pendingLimitStr != "" {
		if val, err := strconv.Atoi(pendingLimitStr); err == nil && val >= 0 {
			pendingLimit = val
		}
	}

	jwtSecret := mustEnv("JWT_SECRET")

	jwtAccessTTLMinutes := 15
	if accessStr := os.Getenv("JWT_ACCESS_TTL_MINUTES"); accessStr != "" {
		if val, err := strconv.Atoi(accessStr); err == nil && val > 0 {
			jwtAccessTTLMinutes = val
		}
	}

	jwtRefreshTTLDays := 7
	if refreshStr := os.Getenv("JWT_REFRESH_TTL_DAYS"); refreshStr != "" {
		if val, err := strconv.Atoi(refreshStr); err == nil && val > 0 {
			jwtRefreshTTLDays = val
		}
	}

	return Config{
		Port:                 port,
		Env:                  env,
		IsProduction:         env == "production",
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		GeminiAPIKey:         os.Getenv("GEMINI_API_KEY"),
		AdminAPIKey:          os.Getenv("ADMIN_API_KEY"),
		CronFetchSchedule:    cronFetch,
		CronProcessSchedule:  cronProcess,
		FrontendURL:          frontendURL,
		AIAPIKey:             os.Getenv("AI_API_KEY"),
		AIBaseURL:            aiBaseURL,
		AIModel:              aiModel,
		PendingLimitPerTopic: pendingLimit,
		JWTSecret:            jwtSecret,
		JWTAccessTTLMinutes:  jwtAccessTTLMinutes,
		JWTRefreshTTLDays:    jwtRefreshTTLDays,
	}
}
