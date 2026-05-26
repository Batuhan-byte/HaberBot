package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	rows, err := conn.Query(ctx, "SELECT id, title, turkish_title, processed_at FROM articles ORDER BY fetched_at DESC LIMIT 10")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("--- Recent Articles in DB ---")
	for rows.Next() {
		var id, title, titleTr string
		var processedAt interface{}
		err := rows.Scan(&id, &title, &titleTr, &processedAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %s\n- Original Title: %s\n- Turkish Title: %s\n- Processed At: %v\n\n", id, title, titleTr, processedAt)
	}
}
