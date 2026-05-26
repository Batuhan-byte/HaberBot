package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/valueobject"
)

// PostgresTopicRepo implements port.TopicRepository using PostgreSQL.
type PostgresTopicRepo struct {
	pool *pgxpool.Pool
}

// NewPostgresTopicRepo creates a new PostgresTopicRepo backed by the given
// connection pool.
func NewPostgresTopicRepo(pool *pgxpool.Pool) *PostgresTopicRepo {
	return &PostgresTopicRepo{pool: pool}
}

// Save persists a topic with an upsert strategy. JSONB columns are used for
// keywords, sources, and rss_feeds arrays.
func (r *PostgresTopicRepo) Save(ctx context.Context, topic *entity.Topic) error {
	keywordsJSON, err := r.marshalKeywords(topic.Keywords)
	if err != nil {
		return fmt.Errorf("marshaling keywords: %w", err)
	}

	sourcesJSON, err := r.marshalSources(topic.Sources)
	if err != nil {
		return fmt.Errorf("marshaling sources: %w", err)
	}

	feedsJSON, err := json.Marshal(topic.RSSFeeds)
	if err != nil {
		return fmt.Errorf("marshaling rss feeds: %w", err)
	}

	query := `
		INSERT INTO topics (id, name, slug, keywords, sources, rss_feeds, is_active, created_at)
		VALUES (COALESCE(NULLIF($1, ''), gen_random_uuid()::text),
			$2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			slug = EXCLUDED.slug,
			keywords = EXCLUDED.keywords,
			sources = EXCLUDED.sources,
			rss_feeds = EXCLUDED.rss_feeds,
			is_active = EXCLUDED.is_active
		RETURNING id, created_at`

	now := time.Now()
	return r.pool.QueryRow(ctx, query,
		topic.ID, topic.Name, topic.Slug,
		keywordsJSON, sourcesJSON, feedsJSON,
		topic.IsActive, now,
	).Scan(&topic.ID, &topic.CreatedAt)
}

// FindByID retrieves a single topic by its unique identifier.
func (r *PostgresTopicRepo) FindByID(ctx context.Context, id string) (*entity.Topic, error) {
	query := `
		SELECT id, name, slug, keywords, sources, rss_feeds, is_active, created_at
		FROM topics WHERE id = $1`

	topic, err := r.scanTopic(r.pool.QueryRow(ctx, query, id))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return topic, err
}

// FindBySlug retrieves a single topic by its URL-friendly slug.
func (r *PostgresTopicRepo) FindBySlug(ctx context.Context, slug string) (*entity.Topic, error) {
	query := `
		SELECT id, name, slug, keywords, sources, rss_feeds, is_active, created_at
		FROM topics WHERE slug = $1`

	topic, err := r.scanTopic(r.pool.QueryRow(ctx, query, slug))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return topic, err
}

// FindAll retrieves all topics.
func (r *PostgresTopicRepo) FindAll(ctx context.Context) ([]*entity.Topic, error) {
	query := `
		SELECT id, name, slug, keywords, sources, rss_feeds, is_active, created_at
		FROM topics ORDER BY created_at ASC`

	return r.queryTopics(ctx, query)
}

// FindActive retrieves only topics that are currently active.
func (r *PostgresTopicRepo) FindActive(ctx context.Context) ([]*entity.Topic, error) {
	query := `
		SELECT id, name, slug, keywords, sources, rss_feeds, is_active, created_at
		FROM topics WHERE is_active = true ORDER BY created_at ASC`

	return r.queryTopics(ctx, query)
}

// Delete removes a topic by its unique identifier.
func (r *PostgresTopicRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM topics WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting topic %s: %w", id, err)
	}
	return nil
}

// queryTopics executes a query and returns a slice of topics.
func (r *PostgresTopicRepo) queryTopics(ctx context.Context, query string, args ...any) ([]*entity.Topic, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying topics: %w", err)
	}
	defer rows.Close()

	var topics []*entity.Topic
	for rows.Next() {
		topic, err := r.scanTopicFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning topic row: %w", err)
		}
		topics = append(topics, topic)
	}

	return topics, rows.Err()
}

// scanTopic scans a single topic from a pgx.Row.
func (r *PostgresTopicRepo) scanTopic(row pgx.Row) (*entity.Topic, error) {
	var topic entity.Topic
	var keywordsJSON, sourcesJSON, feedsJSON []byte

	err := row.Scan(
		&topic.ID, &topic.Name, &topic.Slug,
		&keywordsJSON, &sourcesJSON, &feedsJSON,
		&topic.IsActive, &topic.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return r.unmarshalTopicFields(&topic, keywordsJSON, sourcesJSON, feedsJSON)
}

// scanTopicFromRows scans a single topic from pgx.Rows.
func (r *PostgresTopicRepo) scanTopicFromRows(rows pgx.Rows) (*entity.Topic, error) {
	var topic entity.Topic
	var keywordsJSON, sourcesJSON, feedsJSON []byte

	err := rows.Scan(
		&topic.ID, &topic.Name, &topic.Slug,
		&keywordsJSON, &sourcesJSON, &feedsJSON,
		&topic.IsActive, &topic.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return r.unmarshalTopicFields(&topic, keywordsJSON, sourcesJSON, feedsJSON)
}

// unmarshalTopicFields deserializes JSONB columns into the topic entity.
func (r *PostgresTopicRepo) unmarshalTopicFields(
	topic *entity.Topic,
	keywordsJSON, sourcesJSON, feedsJSON []byte,
) (*entity.Topic, error) {
	var rawKeywords []string
	if err := json.Unmarshal(keywordsJSON, &rawKeywords); err != nil {
		return nil, fmt.Errorf("unmarshaling keywords: %w", err)
	}
	for _, kw := range rawKeywords {
		topic.Keywords = append(topic.Keywords, valueobject.TopicKeyword(kw))
	}

	var rawSources []string
	if err := json.Unmarshal(sourcesJSON, &rawSources); err != nil {
		return nil, fmt.Errorf("unmarshaling sources: %w", err)
	}
	for _, src := range rawSources {
		topic.Sources = append(topic.Sources, valueobject.SourceType(src))
	}

	if err := json.Unmarshal(feedsJSON, &topic.RSSFeeds); err != nil {
		return nil, fmt.Errorf("unmarshaling rss feeds: %w", err)
	}

	return topic, nil
}

// marshalKeywords serializes TopicKeyword slice to JSON bytes.
func (r *PostgresTopicRepo) marshalKeywords(keywords []valueobject.TopicKeyword) ([]byte, error) {
	raw := make([]string, len(keywords))
	for i, kw := range keywords {
		raw[i] = kw.String()
	}
	return json.Marshal(raw)
}

// marshalSources serializes SourceType slice to JSON bytes.
func (r *PostgresTopicRepo) marshalSources(sources []valueobject.SourceType) ([]byte, error) {
	raw := make([]string, len(sources))
	for i, src := range sources {
		raw[i] = src.String()
	}
	return json.Marshal(raw)
}
