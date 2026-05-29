-- Migration 012 Down: Remove security hardening constraints
ALTER TABLE user_favorites DROP CONSTRAINT IF EXISTS chk_no_self_favorite;
ALTER TABLE profile_reports DROP CONSTRAINT IF EXISTS chk_profile_reports_status;
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_users_role;
