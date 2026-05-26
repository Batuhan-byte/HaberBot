package main

import (
	"context"
	"log"
	"os"
	"time"

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

	// Get the topic ID for "Yapay Zeka"
	var topicID string
	err = conn.QueryRow(ctx, "SELECT id FROM topics WHERE slug = 'yapay-zeka'").Scan(&topicID)
	if err != nil {
		log.Fatal("Could not find Yapay Zeka topic: ", err)
	}

	// Insert a sample article
	articleID := uuid.New().String()
	_, err = conn.Exec(ctx, `
		INSERT INTO articles (
			id, title, turkish_title, original_url, source_type, 
			original_content, turkish_summary, score, topic_id, 
			image_url, processed_at, fetched_at, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
		ON CONFLICT (original_url) DO NOTHING`,
		articleID,
		"Google DeepMind introduces Gemini 1.5 Pro",
		"Google DeepMind, Gemini 1.5 Pro'yu tanıttı",
		"https://deepmind.google/technologies/gemini/1-5-pro-sample",
		"hackernews",
		"Google DeepMind introduces Gemini 1.5 Pro with a context window of up to 1 million tokens...",
		"Google DeepMind, yapay zeka dünyasında devrim yaratan Gemini 1.5 Pro modelini duyurdu. Yeni sürüm, 1 milyon tokenlık muazzam bağlam kapasitesi sayesinde binlerce sayfalık metni ve saatlerce süren videoları tek seferde analiz edebiliyor.",
		2540,
		topicID,
		"https://upload.wikimedia.org/wikipedia/commons/8/8a/Google_Gemini_logo.svg",
		time.Now(),
		time.Now().Add(-24*time.Hour),
		time.Now().Add(-24*time.Hour),
	)

	if err != nil {
		log.Fatal("Failed to insert article: ", err)
	}
	log.Println("Successfully added a sample article to the database!")
}
