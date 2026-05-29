package entity

import "time"

// Comment represents a user comment on an article.
type Comment struct {
	ID        string    `json:"id"`
	ArticleID string    `json:"article_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username,omitempty"` // Populated during joins
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
