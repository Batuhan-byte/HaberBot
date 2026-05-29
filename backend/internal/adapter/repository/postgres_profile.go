package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"haberbot/internal/domain/entity"
)

// PostgresProfileRepo implements port.ProfileRepository using PostgreSQL.
type PostgresProfileRepo struct {
	pool *pgxpool.Pool
}

// NewPostgresProfileRepo creates a new PostgresProfileRepo.
func NewPostgresProfileRepo(pool *pgxpool.Pool) *PostgresProfileRepo {
	return &PostgresProfileRepo{pool: pool}
}

// GetStatsByUserID retrieves stats for a given user ID.
func (r *PostgresProfileRepo) GetStatsByUserID(ctx context.Context, userID string) (*entity.UserStats, error) {
	query := `SELECT user_id, articles_published, updated_at FROM user_stats WHERE user_id = $1`
	var stats entity.UserStats
	err := r.pool.QueryRow(ctx, query, userID).Scan(&stats.UserID, &stats.ArticlesPublished, &stats.UpdatedAt)
	if err == pgx.ErrNoRows {
		// Return a default zero stats object instead of error
		return &entity.UserStats{UserID: userID, ArticlesPublished: 0, UpdatedAt: time.Now()}, nil
	}
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// UpsertStats inserts or updates user statistics.
func (r *PostgresProfileRepo) UpsertStats(ctx context.Context, stats *entity.UserStats) error {
	query := `
		INSERT INTO user_stats (user_id, articles_published, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			articles_published = EXCLUDED.articles_published,
			updated_at = NOW()`
	_, err := r.pool.Exec(ctx, query, stats.UserID, stats.ArticlesPublished)
	return err
}

// AddFavorite records a new favorite relationship.
func (r *PostgresProfileRepo) AddFavorite(ctx context.Context, favorite *entity.UserFavorite) error {
	query := `
		INSERT INTO user_favorites (id, user_id, favorite_user_id, created_at)
		VALUES (COALESCE(NULLIF($1, ''), gen_random_uuid()::text), $2, $3, NOW())
		RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, favorite.ID, favorite.UserID, favorite.FavoriteUserID).
		Scan(&favorite.ID, &favorite.CreatedAt)
}

// RemoveFavorite deletes a favorite relationship.
func (r *PostgresProfileRepo) RemoveFavorite(ctx context.Context, userID, favoriteUserID string) error {
	query := `DELETE FROM user_favorites WHERE user_id = $1 AND favorite_user_id = $2`
	result, err := r.pool.Exec(ctx, query, userID, favoriteUserID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// IsFavorite checks if a favorite relationship exists.
func (r *PostgresProfileRepo) IsFavorite(ctx context.Context, userID, favoriteUserID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_favorites WHERE user_id = $1 AND favorite_user_id = $2)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, userID, favoriteUserID).Scan(&exists)
	return exists, err
}

// ListFavorites retrieves the list of users favorited by a user.
func (r *PostgresProfileRepo) ListFavorites(ctx context.Context, userID string, limit, offset int) ([]*entity.User, int, error) {
	countQuery := `SELECT COUNT(*) FROM user_favorites WHERE user_id = $1`
	var total int
	err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT u.id, u.username, u.email, u.role, u.bio, u.avatar_url, u.join_date, u.created_at, u.updated_at
		FROM user_favorites f
		JOIN users u ON f.favorite_user_id = u.id
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var favorites []*entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Role, &user.Bio, &user.AvatarURL, &user.JoinDate, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		favorites = append(favorites, &user)
	}

	return favorites, total, nil
}

// SubmitReport files a new profile report.
func (r *PostgresProfileRepo) SubmitReport(ctx context.Context, report *entity.ProfileReport) error {
	query := `
		INSERT INTO profile_reports (id, reporter_user_id, reported_user_id, reason, comment, status, created_at)
		VALUES (COALESCE(NULLIF($1, ''), gen_random_uuid()::text), $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, report.ID, report.ReporterUserID, report.ReportedUserID, report.Reason, report.Comment, report.Status).
		Scan(&report.ID, &report.CreatedAt)
}

// HasRecentReport checks if a report has been submitted in the last 24 hours.
func (r *PostgresProfileRepo) HasRecentReport(ctx context.Context, reporterUserID, reportedUserID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM profile_reports 
			WHERE reporter_user_id = $1 AND reported_user_id = $2 AND created_at > NOW() - INTERVAL '24 hours'
		)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, reporterUserID, reportedUserID).Scan(&exists)
	return exists, err
}

// RecountStats counts total comments posted by a user and upserts it into user_stats table.
func (r *PostgresProfileRepo) RecountStats(ctx context.Context, userID string) error {
	var count int
	countQuery := `SELECT COUNT(*) FROM comments WHERE user_id = $1`
	err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&count)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO user_stats (user_id, articles_published, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			articles_published = EXCLUDED.articles_published,
			updated_at = NOW()`
	_, err = r.pool.Exec(ctx, query, userID, count)
	return err
}
