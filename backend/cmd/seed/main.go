package main

import (
	"context"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	// Seed AI topic
	id := uuid.New().String()
	_, err = conn.Exec(ctx, `
		INSERT INTO topics (id, name, slug, keywords, sources, rss_feeds, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (slug) DO UPDATE SET
			keywords = EXCLUDED.keywords,
			sources = EXCLUDED.sources,
			rss_feeds = EXCLUDED.rss_feeds`,
		id, "Yapay Zeka", "yapay-zeka", `["AI", "LLM", "Machine Learning", "Google", "OpenAI", "Tech", "Developer", "Software", "Python", "Rust", "Microsoft"]`, `["hackernews", "rss"]`, `["https://news.ycombinator.com/rss", "https://dev.to/feed", "https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml"]`, true)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Seeded Yapay Zeka topic")
}
