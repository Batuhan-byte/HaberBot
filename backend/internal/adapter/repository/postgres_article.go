// Package repository provides concrete database implementations of the
// domain port interfaces using PostgreSQL via pgx.
package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/valueobject"
)

// PostgresArticleRepo implements port.ArticleRepository using PostgreSQL.
type PostgresArticleRepo struct {
	pool *pgxpool.Pool
}

// NewPostgresArticleRepo creates a new PostgresArticleRepo backed by the given
// connection pool.
func NewPostgresArticleRepo(pool *pgxpool.Pool) *PostgresArticleRepo {
	return &PostgresArticleRepo{pool: pool}
}

// Save persists an article. If the article has an ID, it performs an upsert;
// otherwise it inserts a new record with a generated UUID.
func (r *PostgresArticleRepo) Save(ctx context.Context, article *entity.Article) error {
	query := `
		INSERT INTO articles (id, title, turkish_title, original_url, source_type,
			original_content, turkish_content, turkish_summary, score, topic_id, image_url,
			processed_at, fetched_at, created_at, is_approved, is_hidden, approved_at)
		VALUES (COALESCE(NULLIF($1, ''), gen_random_uuid()::text),
			$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		ON CONFLICT (id) DO UPDATE SET
			turkish_title = EXCLUDED.turkish_title,
			turkish_content = EXCLUDED.turkish_content,
			turkish_summary = EXCLUDED.turkish_summary,
			processed_at = EXCLUDED.processed_at,
			score = EXCLUDED.score,
			is_approved = EXCLUDED.is_approved,
			is_hidden = EXCLUDED.is_hidden,
			approved_at = EXCLUDED.approved_at
		RETURNING id, created_at`

	article.Title = strings.ToValidUTF8(article.Title, "")
	article.TurkishTitle = strings.ToValidUTF8(article.TurkishTitle, "")
	article.OriginalContent = strings.ToValidUTF8(article.OriginalContent, "")
	article.TurkishContent = strings.ToValidUTF8(article.TurkishContent, "")
	article.TurkishSummary = strings.ToValidUTF8(article.TurkishSummary, "")

	now := time.Now()
	if article.FetchedAt.IsZero() {
		article.FetchedAt = now
	}

	var approvedAt *time.Time
	if article.IsApproved {
		if article.ApprovedAt != nil {
			approvedAt = article.ApprovedAt
		} else {
			approvedAt = &now
		}
	}

	return r.pool.QueryRow(ctx, query,
		article.ID, article.Title, article.TurkishTitle, article.OriginalURL,
		article.SourceType.String(), article.OriginalContent, article.TurkishContent, article.TurkishSummary,
		article.Score, article.TopicID, article.ImageURL, article.ProcessedAt,
		article.FetchedAt, now, article.IsApproved, article.IsHidden, approvedAt,
	).Scan(&article.ID, &article.CreatedAt)
}

// FindByID retrieves a single article by its unique identifier.
func (r *PostgresArticleRepo) FindByID(ctx context.Context, id string) (*entity.Article, error) {
	query := `
			SELECT id, title, turkish_title, original_url, source_type,
				original_content, turkish_content, turkish_summary, score, topic_id, image_url,
				processed_at, fetched_at, created_at, is_approved, is_hidden, approved_at
			FROM articles WHERE id = $1`

	article, err := r.scanArticle(r.pool.QueryRow(ctx, query, id))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return article, err
}

// FindByTopicID retrieves paginated articles for a given topic (visitors only).
func (r *PostgresArticleRepo) FindByTopicID(ctx context.Context, topicID string, limit, offset int) ([]*entity.Article, error) {
	query := `
			SELECT id, title, turkish_title, original_url, source_type,
				original_content, turkish_content, turkish_summary, score, topic_id, image_url,
				processed_at, fetched_at, created_at, is_approved, is_hidden, approved_at
			FROM articles
		WHERE topic_id = $1 AND is_approved = true AND is_hidden = false
		ORDER BY fetched_at DESC
		LIMIT $2 OFFSET $3`

	return r.queryArticles(ctx, query, topicID, limit, offset)
}

// FindRecent retrieves the most recently fetched processed articles (visitors only).
func (r *PostgresArticleRepo) FindRecent(ctx context.Context, limit int) ([]*entity.Article, error) {
	query := `
			SELECT id, title, turkish_title, original_url, source_type,
				original_content, turkish_content, turkish_summary, score, topic_id, image_url,
				processed_at, fetched_at, created_at, is_approved, is_hidden, approved_at
			FROM articles
		WHERE is_approved = true AND is_hidden = false
		ORDER BY fetched_at DESC
		LIMIT $1`

	return r.queryArticles(ctx, query, limit)
}

// FindUnprocessed retrieves articles that have not been processed by AI yet.
func (r *PostgresArticleRepo) FindUnprocessed(ctx context.Context, limit int) ([]*entity.Article, error) {
	query := `
			SELECT id, title, turkish_title, original_url, source_type,
				original_content, turkish_content, turkish_summary, score, topic_id, image_url,
				processed_at, fetched_at, created_at, is_approved, is_hidden, approved_at
			FROM articles
		WHERE processed_at IS NULL
		ORDER BY fetched_at ASC
		LIMIT $1`

	return r.queryArticles(ctx, query, limit)
}

// ExistsByURL checks whether an article with the given URL already exists.
func (r *PostgresArticleRepo) ExistsByURL(ctx context.Context, url string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM articles WHERE original_url = $1)`
	err := r.pool.QueryRow(ctx, query, url).Scan(&exists)
	return exists, err
}

// CountByTopicID returns the total number of approved/visible articles for a topic.
func (r *PostgresArticleRepo) CountByTopicID(ctx context.Context, topicID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM articles WHERE topic_id = $1 AND is_approved = true AND is_hidden = false`
	err := r.pool.QueryRow(ctx, query, topicID).Scan(&count)
	return count, err
}

// Search searches for approved and visible articles matching the query.
func (r *PostgresArticleRepo) Search(ctx context.Context, query string, source valueobject.SourceType, limit, offset int) ([]*entity.Article, error) {
	var queryStr string
	var args []interface{}
	
	if source == "" {
		searchTerm := "%" + query + "%"
		queryStr = `
			SELECT id, title, turkish_title, original_url, source_type,
				original_content, turkish_content, turkish_summary, score, topic_id, image_url,
				processed_at, fetched_at, created_at, is_approved, is_hidden, approved_at
			FROM articles
			WHERE (title ILIKE $1 OR turkish_title ILIKE $1 OR turkish_summary ILIKE $1 OR original_content ILIKE $1)
			AND is_approved = true AND is_hidden = false
			ORDER BY 
				CASE 
					WHEN title ILIKE $1 THEN 3
					WHEN turkish_title ILIKE $1 THEN 2
					WHEN turkish_summary ILIKE $1 THEN 1
					ELSE 0
				END DESC,
				fetched_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{searchTerm, limit, offset}
	} else {
		searchTerm := "%" + query + "%"
		queryStr = `
			SELECT id, title, turkish_title, original_url, source_type,
				original_content, turkish_content, turkish_summary, score, topic_id, image_url,
				processed_at, fetched_at, created_at, is_approved, is_hidden, approved_at
			FROM articles
			WHERE (title ILIKE $1 OR turkish_title ILIKE $1 OR turkish_summary ILIKE $1 OR original_content ILIKE $1)
			AND source_type = $2 AND is_approved = true AND is_hidden = false
			ORDER BY 
				CASE 
					WHEN title ILIKE $1 THEN 3
					WHEN turkish_title ILIKE $1 THEN 2
					WHEN turkish_summary ILIKE $1 THEN 1
					ELSE 0
				END DESC,
				fetched_at DESC
			LIMIT $3 OFFSET $4
		`
		args = []interface{}{searchTerm, source.String(), limit, offset}
	}

	return r.queryArticles(ctx, queryStr, args...)
}

// UpdateApprovalStatus updates the approval status of an article.
func (r *PostgresArticleRepo) UpdateApprovalStatus(ctx context.Context, id string, isApproved bool) error {
	var query string
	if isApproved {
		query = `UPDATE articles SET is_approved = $1, approved_at = NOW() WHERE id = $2`
	} else {
		query = `UPDATE articles SET is_approved = $1, approved_at = NULL WHERE id = $2`
	}
	_, err := r.pool.Exec(ctx, query, isApproved, id)
	return err
}

// UpdateHidingStatus updates the hiding status of an article.
func (r *PostgresArticleRepo) UpdateHidingStatus(ctx context.Context, id string, isHidden bool) error {
	query := `UPDATE articles SET is_hidden = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, isHidden, id)
	return err
}

// Delete hard-deletes an article from the database.
func (r *PostgresArticleRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM articles WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// FindAllAdmin retrieves all articles for the admin panel with optional topic filtering.
func (r *PostgresArticleRepo) FindAllAdmin(ctx context.Context, topicID string, limit, offset int) ([]*entity.Article, error) {
	var query string
	var args []any

	if topicID == "" {
		query = `
			SELECT id, title, turkish_title, original_url, source_type,
				original_content, turkish_content, turkish_summary, score, topic_id, image_url,
				processed_at, fetched_at, created_at, is_approved, is_hidden, approved_at
			FROM articles
			WHERE is_approved = true
			ORDER BY fetched_at DESC
			LIMIT $1 OFFSET $2`
		args = []any{limit, offset}
	} else if topicID == "pending" || strings.HasPrefix(topicID, "pending:") {
		var actualTopicID string
		if strings.HasPrefix(topicID, "pending:") {
			actualTopicID = strings.TrimPrefix(topicID, "pending:")
		}

		if actualTopicID == "" {
			query = `
				SELECT id, title, turkish_title, original_url, source_type,
					original_content, turkish_content, turkish_summary, score, topic_id, image_url,
					processed_at, fetched_at, created_at, is_approved, is_hidden, approved_at
				FROM articles
				WHERE is_approved = false
				ORDER BY fetched_at DESC
				LIMIT $1 OFFSET $2`
			args = []any{limit, offset}
		} else {
			query = `
				SELECT id, title, turkish_title, original_url, source_type,
					original_content, turkish_content, turkish_summary, score, topic_id, image_url,
					processed_at, fetched_at, created_at, is_approved, is_hidden, approved_at
				FROM articles
				WHERE is_approved = false AND topic_id = $1
				ORDER BY fetched_at DESC
				LIMIT $2 OFFSET $3`
			args = []any{actualTopicID, limit, offset}
		}
	} else {
		query = `
			SELECT id, title, turkish_title, original_url, source_type,
				original_content, turkish_content, turkish_summary, score, topic_id, image_url,
				processed_at, fetched_at, created_at, is_approved, is_hidden, approved_at
			FROM articles
			WHERE topic_id = $1 AND is_approved = true
			ORDER BY fetched_at DESC
			LIMIT $2 OFFSET $3`
		args = []any{topicID, limit, offset}
	}

	return r.queryArticles(ctx, query, args...)
}

// CountAllAdmin returns the total count of articles for the admin panel.
func (r *PostgresArticleRepo) CountAllAdmin(ctx context.Context, topicID string) (int, error) {
	var query string
	var args []any
	var count int

	if topicID == "" {
		query = `SELECT COUNT(*) FROM articles WHERE is_approved = true`
		args = []any{}
	} else if topicID == "pending" || strings.HasPrefix(topicID, "pending:") {
		var actualTopicID string
		if strings.HasPrefix(topicID, "pending:") {
			actualTopicID = strings.TrimPrefix(topicID, "pending:")
		}

		if actualTopicID == "" {
			query = `SELECT COUNT(*) FROM articles WHERE is_approved = false`
			args = []any{}
		} else {
			query = `SELECT COUNT(*) FROM articles WHERE is_approved = false AND topic_id = $1`
			args = []any{actualTopicID}
		}
	} else {
		query = `SELECT COUNT(*) FROM articles WHERE topic_id = $1 AND is_approved = true`
		args = []any{topicID}
	}

	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// queryArticles executes a query and returns a slice of articles.
func (r *PostgresArticleRepo) queryArticles(ctx context.Context, query string, args ...any) ([]*entity.Article, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying articles: %w", err)
	}
	defer rows.Close()

	var articles []*entity.Article
	for rows.Next() {
		article, err := r.scanArticleFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning article row: %w", err)
		}
		articles = append(articles, article)
	}

	return articles, rows.Err()
}

// scanArticle scans a single article from a pgx.Row.
func (r *PostgresArticleRepo) scanArticle(row pgx.Row) (*entity.Article, error) {
	var article entity.Article
	var sourceType string

	err := row.Scan(
		&article.ID, &article.Title, &article.TurkishTitle, &article.OriginalURL,
		&sourceType, &article.OriginalContent, &article.TurkishContent, &article.TurkishSummary,
		&article.Score, &article.TopicID, &article.ImageURL,
		&article.ProcessedAt, &article.FetchedAt, &article.CreatedAt,
		&article.IsApproved, &article.IsHidden, &article.ApprovedAt,
	)
	if err != nil {
		return nil, err
	}

	article.SourceType = valueobject.SourceType(sourceType)
	return &article, nil
}

// scanArticleFromRows scans a single article from pgx.Rows.
func (r *PostgresArticleRepo) scanArticleFromRows(rows pgx.Rows) (*entity.Article, error) {
	var article entity.Article
	var sourceType string

	err := rows.Scan(
		&article.ID, &article.Title, &article.TurkishTitle, &article.OriginalURL,
		&sourceType, &article.OriginalContent, &article.TurkishContent, &article.TurkishSummary,
		&article.Score, &article.TopicID, &article.ImageURL,
		&article.ProcessedAt, &article.FetchedAt, &article.CreatedAt,
		&article.IsApproved, &article.IsHidden, &article.ApprovedAt,
	)
	if err != nil {
		return nil, err
	}

	article.SourceType = valueobject.SourceType(sourceType)
	return &article, nil
}

// TrimPendingByTopic removes oldest pending articles of a given topic (or NULL topic) if count exceeds the limit.
func (r *PostgresArticleRepo) TrimPendingByTopic(ctx context.Context, topicID string, limit int) error {
	var query string
	var args []any

	if topicID == "" {
		query = `
			DELETE FROM articles
			WHERE topic_id IS NULL AND is_approved = false AND id IN (
				SELECT id FROM articles
				WHERE topic_id IS NULL AND is_approved = false
				ORDER BY fetched_at DESC
				OFFSET $1
			)`
		args = []any{limit}
	} else {
		query = `
			DELETE FROM articles
			WHERE topic_id = $1 AND is_approved = false AND id IN (
				SELECT id FROM articles
				WHERE topic_id = $1 AND is_approved = false
				ORDER BY fetched_at DESC
				OFFSET $2
			)`
		args = []any{topicID, limit}
	}

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}
