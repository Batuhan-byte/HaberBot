package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"haberbot/internal/infrastructure/config"
	"haberbot/internal/infrastructure/container"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	cfg := config.Load()
	c, err := container.NewContainer(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer c.DBPool.Close()

	log.Println("Starting fast global translation of ALL unprocessed articles in database using Mistral...")

	totalProcessed := 0
	for {
		// Fetch 5 unprocessed articles at a time
		articles, err := c.ArticleRepo.FindUnprocessed(ctx, 5)
		if err != nil {
			log.Fatalf("Error finding unprocessed: %v", err)
		}

		if len(articles) == 0 {
			log.Println("All articles processed and translated successfully!")
			break
		}

		log.Printf("Found %d unprocessed articles. Processing batch...", len(articles))

		for _, article := range articles {
			log.Printf("Translating: %s", article.Title)
			
			turkishTitle, turkishContent, err := c.AIProcessor.Translate(ctx, article.Title, article.OriginalContent)
			if err != nil {
				log.Printf("Error translating article %s: %v", article.ID, err)
				time.Sleep(2 * time.Second)
				continue
			}

			now := time.Now()
			article.TurkishTitle = turkishTitle
			article.TurkishContent = turkishContent
			article.ProcessedAt = &now

			if err := c.ArticleRepo.Save(ctx, article); err != nil {
				log.Fatalf("Error saving article %s: %v", article.ID, err)
			}

			totalProcessed++
			log.Printf("Successfully translated to: %s", turkishTitle)

			// Since we use Mistral AI, we only need a tiny 1.5s delay to be perfectly safe under rate limits
			time.Sleep(1500 * time.Millisecond)
		}
	}

	log.Printf("Fast global translation complete! Total translated: %d", totalProcessed)
}
