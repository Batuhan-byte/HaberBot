package entity

import "time"

// UserStats represents the cached publication statistics for a user.
type UserStats struct {
	UserID            string    `json:"user_id"`
	ArticlesPublished int       `json:"articles_published"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// UserFavorite represents a user's favorited profile connection.
type UserFavorite struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	FavoriteUserID string    `json:"favorite_user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// ProfileReport represents a user profile report for moderation.
type ProfileReport struct {
	ID             string     `json:"id"`
	ReporterUserID string     `json:"reporter_user_id"`
	ReportedUserID string     `json:"reported_user_id"`
	Reason         string     `json:"reason"`
	Comment        string     `json:"comment,omitempty"`
	Status         string     `json:"status"` // "open", "reviewed", "dismissed", "resolved"
	CreatedAt      time.Time  `json:"created_at"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
}
