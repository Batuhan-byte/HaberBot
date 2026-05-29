package entity

import "time"

// User represents an authenticated user in the system.
type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         string     `json:"role"` // "Admin" or "User"
	RefreshToken string     `json:"-"`
	Bio          *string    `json:"bio"`
	AvatarURL    *string    `json:"avatar_url"`
	JoinDate     time.Time  `json:"join_date"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
