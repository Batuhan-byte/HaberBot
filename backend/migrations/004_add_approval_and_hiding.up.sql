-- Migration Up: Add is_approved and is_hidden columns to articles table
ALTER TABLE articles ADD COLUMN is_approved BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE articles ADD COLUMN is_hidden BOOLEAN NOT NULL DEFAULT false;

-- Approve all existing articles so they remain public
UPDATE articles SET is_approved = true;
