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

	log.Println("All migrations executed successfully.")
}
