-- Migration Up: Make topic_id column nullable in articles table
ALTER TABLE articles ALTER COLUMN topic_id DROP NOT NULL;
