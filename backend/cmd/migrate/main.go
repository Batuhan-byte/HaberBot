package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

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

	// Read and execute 001
	sql1, err := os.ReadFile("migrations/001_create_topics.up.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Exec(ctx, string(sql1))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Executed 001_create_topics.up.sql")

	// Read and execute 002
	sql2, err := os.ReadFile("migrations/002_create_articles.up.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Exec(ctx, string(sql2))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Executed 002_create_articles.up.sql")

	// Read and execute 003
	sql3, err := os.ReadFile("migrations/003_add_turkish_content.up.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Exec(ctx, string(sql3))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Executed 003_add_turkish_content.up.sql")

	// Read and execute 004
	sql4, err := os.ReadFile("migrations/004_add_approval_and_hiding.up.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Exec(ctx, string(sql4))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Executed 004_add_approval_and_hiding.up.sql")

	// Read and execute 005
	sql5, err := os.ReadFile("migrations/005_add_approved_at.up.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Exec(ctx, string(sql5))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Executed 005_add_approved_at.up.sql")

	// Read and execute 006
	sql6, err := os.ReadFile("migrations/006_make_topic_id_nullable.up.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Exec(ctx, string(sql6))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Executed 006_make_topic_id_nullable.up.sql")

	// Read and execute 007
	sql7, err := os.ReadFile("migrations/007_create_users.up.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Exec(ctx, string(sql7))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Executed 007_create_users.up.sql")

	// Read and execute 008
	sql8, err := os.ReadFile("migrations/008_create_comments.up.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Exec(ctx, string(sql8))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Executed 008_create_comments.up.sql")

	// Read and execute 009
	sql9, err := os.ReadFile("migrations/009_seed_admin.up.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Exec(ctx, string(sql9))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Executed 009_seed_admin.up.sql")

	log.Println("All migrations executed successfully.")
}
