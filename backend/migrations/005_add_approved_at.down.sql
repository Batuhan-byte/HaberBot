-- Migration Down: Remove approved_at column from articles table
ALTER TABLE articles DROP COLUMN approved_at;
