-- Migration Down: Remove is_approved and is_hidden columns from articles table
ALTER TABLE articles DROP COLUMN is_approved;
ALTER TABLE articles DROP COLUMN is_hidden;
