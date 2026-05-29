-- Migration Up: Add approved_at column to articles table
ALTER TABLE articles ADD COLUMN approved_at TIMESTAMP WITH TIME ZONE;

-- For existing approved articles, set approved_at to created_at
UPDATE articles SET approved_at = created_at WHERE is_approved = true;
