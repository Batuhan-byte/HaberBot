package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func execSQL(ctx context.Context, conn *pgx.Conn, filepath string) {
	sql, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read migration file %s: %v", filepath, err)
	}

	_, err = conn.Exec(ctx, string(sql))
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "already exists") || strings.Contains(errStr, "duplicate key") || strings.Contains(errStr, "already a member") {
			log.Printf("[Skipped] %s: already executed (resource exists)\n", filepath)
			return
		}
		log.Fatalf("Failed to execute migration %s: %v", filepath, err)
	}
	log.Printf("[Executed] %s\n", filepath)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	// List of all migrations in order
	migrations := []string{
		"migrations/001_create_topics.up.sql",
		"migrations/002_create_articles.up.sql",
		"migrations/003_add_turkish_content.up.sql",
		"migrations/004_add_approval_and_hiding.up.sql",
		"migrations/005_add_approved_at.up.sql",
		"migrations/006_make_topic_id_nullable.up.sql",
		"migrations/007_create_users.up.sql",
		"migrations/008_create_comments.up.sql",
		"migrations/009_seed_admin.up.sql",
		"migrations/010_add_email_to_users.up.sql",
		"migrations/011_create_user_profiles.up.sql",
		"migrations/012_security_constraints.up.sql",
	}

	for _, m := range migrations {
		execSQL(ctx, conn, m)
	}

	log.Println("All migrations processed successfully.")
}
