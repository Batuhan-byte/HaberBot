package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"haberbot/internal/domain/entity"
)

// PostgresUserRepo implements port.UserRepository using PostgreSQL.
type PostgresUserRepo struct {
	pool *pgxpool.Pool
}

// NewPostgresUserRepo creates a new PostgresUserRepo.
func NewPostgresUserRepo(pool *pgxpool.Pool) *PostgresUserRepo {
	return &PostgresUserRepo{pool: pool}
}

// Save persists a user (insert or update).
func (r *PostgresUserRepo) Save(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, role, refresh_token, created_at, updated_at)
		VALUES (COALESCE(NULLIF($1, ''), gen_random_uuid()::text), $2, $3, $4, $5, $6, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			username = EXCLUDED.username,
			email = EXCLUDED.email,
			password_hash = EXCLUDED.password_hash,
			role = EXCLUDED.role,
			refresh_token = EXCLUDED.refresh_token,
			updated_at = NOW()
		RETURNING id, created_at, updated_at`

	var refreshToken *string
	if user.RefreshToken != "" {
		refreshToken = &user.RefreshToken
	}

	return r.pool.QueryRow(ctx, query, user.ID, user.Username, user.Email, user.PasswordHash, user.Role, refreshToken).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// FindByID retrieves a user by ID.
func (r *PostgresUserRepo) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT id, username, email, password_hash, role, refresh_token, created_at, updated_at FROM users WHERE id = $1`
	var user entity.User
	var refreshToken *string
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &refreshToken, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if refreshToken != nil {
		user.RefreshToken = *refreshToken
	}
	return &user, nil
}

// FindByUsername retrieves a user by username.
func (r *PostgresUserRepo) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `SELECT id, username, email, password_hash, role, refresh_token, created_at, updated_at FROM users WHERE username = $1`
	var user entity.User
	var refreshToken *string
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &refreshToken, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if refreshToken != nil {
		user.RefreshToken = *refreshToken
	}
	return &user, nil
}

// FindByEmail retrieves a user by email.
func (r *PostgresUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, username, email, password_hash, role, refresh_token, created_at, updated_at FROM users WHERE email = $1`
	var user entity.User
	var refreshToken *string
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &refreshToken, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if refreshToken != nil {
		user.RefreshToken = *refreshToken
	}
	return &user, nil
}

// UpdateRefreshToken updates only the refresh token field of a user.
func (r *PostgresUserRepo) UpdateRefreshToken(ctx context.Context, userID string, token string) error {
	var tokenVal *string
	if token != "" {
		tokenVal = &token
	}
	query := `UPDATE users SET refresh_token = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, tokenVal, userID)
	return err
}
