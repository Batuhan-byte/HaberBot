-- Migration Down: Drop tables and columns added for user profiles

DROP TABLE IF EXISTS profile_reports;
DROP TABLE IF EXISTS user_favorites;
DROP TABLE IF EXISTS user_stats;

ALTER TABLE users DROP COLUMN IF EXISTS join_date;
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
ALTER TABLE users DROP COLUMN IF EXISTS bio;
