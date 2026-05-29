-- Migration Up: Add bio, avatar_url, join_date to users, and create user_stats, user_favorites, profile_reports tables

ALTER TABLE users ADD COLUMN IF NOT EXISTS bio VARCHAR(160) DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(255) DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS join_date TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Backfill join_date with created_at
UPDATE users SET join_date = created_at WHERE join_date IS NULL;

-- 1. user_stats table
CREATE TABLE IF NOT EXISTS user_stats (
    user_id VARCHAR(36) PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    articles_published INT DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Backfill user_stats for existing users
INSERT INTO user_stats (user_id, articles_published, updated_at)
SELECT id, 0, NOW() FROM users
ON CONFLICT (user_id) DO NOTHING;

-- 2. user_favorites table
CREATE TABLE IF NOT EXISTS user_favorites (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    favorite_user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_user_favorite UNIQUE (user_id, favorite_user_id)
);

CREATE INDEX IF NOT EXISTS idx_user_favorites_user_id ON user_favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_user_favorites_favorite_user_id ON user_favorites(favorite_user_id);

-- 3. profile_reports table
CREATE TABLE IF NOT EXISTS profile_reports (
    id VARCHAR(36) PRIMARY KEY,
    reporter_user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reported_user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason VARCHAR(50) NOT NULL DEFAULT 'other',
    comment VARCHAR(500),
    status VARCHAR(50) NOT NULL DEFAULT 'open',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    resolved_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT chk_reporter_reported CHECK (reporter_user_id <> reported_user_id)
);

CREATE INDEX IF NOT EXISTS idx_profile_reports_reported_user_id ON profile_reports(reported_user_id);
