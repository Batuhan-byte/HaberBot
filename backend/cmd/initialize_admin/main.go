// cmd/initialize_admin/main.go
// BUG-003: Secure admin initialization tool.
// This reads ADMIN_USERNAME, ADMIN_PASSWORD, ADMIN_EMAIL from env and creates an admin account.
// Usage: go run cmd/initialize_admin/main.go
// Required env: DATABASE_URL, JWT_SECRET, ADMIN_USERNAME, ADMIN_PASSWORD, ADMIN_EMAIL
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load .env if present
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("FATAL: DATABASE_URL environment variable is required")
	}

	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminEmail := os.Getenv("ADMIN_EMAIL")

	if adminUsername == "" || adminPassword == "" || adminEmail == "" {
		log.Fatal("FATAL: ADMIN_USERNAME, ADMIN_PASSWORD, and ADMIN_EMAIL environment variables are all required")
	}

	if len(adminPassword) < 12 {
		log.Fatal("FATAL: ADMIN_PASSWORD must be at least 12 characters long")
	}

	// Hash the password with bcrypt cost 14 (higher for admin accounts)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), 14)
	if err != nil {
		log.Fatalf("FATAL: failed to hash password: %v", err)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("FATAL: failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	adminID := uuid.New().String()

	_, err = conn.Exec(ctx, `
		INSERT INTO users (id, username, password_hash, role, email)
		VALUES ($1, $2, $3, 'Admin', $4)
		ON CONFLICT (username) DO UPDATE SET
			password_hash = EXCLUDED.password_hash,
			email = EXCLUDED.email,
			role = 'Admin'
	`, adminID, adminUsername, string(hashedPassword), adminEmail)

	if err != nil {
		log.Fatalf("FATAL: failed to create admin user: %v", err)
	}

	fmt.Printf("✅ Admin user %q created/updated successfully.\n", adminUsername)
	fmt.Println("⚠️  IMPORTANT: Clear ADMIN_PASSWORD from your environment immediately after running this command.")
}
