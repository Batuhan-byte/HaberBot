package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"haberbot/internal/domain/entity"
)

// PostgresCommentRepo implements port.CommentRepository using PostgreSQL.
type PostgresCommentRepo struct {
	pool *pgxpool.Pool
}

// NewPostgresCommentRepo creates a new PostgresCommentRepo.
func NewPostgresCommentRepo(pool *pgxpool.Pool) *PostgresCommentRepo {
	return &PostgresCommentRepo{pool: pool}
}

// Save persists a user comment.
func (r *PostgresCommentRepo) Save(ctx context.Context, comment *entity.Comment) error {
	query := `
		INSERT INTO comments (id, article_id, user_id, content, created_at)
		VALUES (COALESCE(NULLIF($1, ''), gen_random_uuid()::text), $2, $3, $4, NOW())
		RETURNING id, created_at`

	return r.pool.QueryRow(ctx, query, comment.ID, comment.ArticleID, comment.UserID, comment.Content).
		Scan(&comment.ID, &comment.CreatedAt)
}

// FindByArticleID retrieves comments belonging to an article, populated with usernames.
func (r *PostgresCommentRepo) FindByArticleID(ctx context.Context, articleID string) ([]*entity.Comment, error) {
	query := `
		SELECT c.id, c.article_id, c.user_id, u.username, c.content, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.article_id = $1
		ORDER BY c.created_at ASC`

	rows, err := r.pool.Query(ctx, query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entity.Comment
	for rows.Next() {
		var c entity.Comment
		err := rows.Scan(&c.ID, &c.ArticleID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, rows.Err()
}
