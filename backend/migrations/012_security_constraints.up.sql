-- Migration 012: Security hardening constraints
-- BUG-006: Add CHECK constraint on users.role column
-- BUG-037: Add CHECK constraint on profile_reports.status
-- BUG-038: Add self-favorite prevention constraint on user_favorites

-- BUG-006: Ensure only valid roles are stored at DB level.
-- This prevents privilege escalation even if application-layer checks are bypassed.
ALTER TABLE users
    ADD CONSTRAINT chk_users_role
    CHECK (role IN ('User', 'Admin'));

-- BUG-037: Ensure profile_reports.status only accepts known values.
ALTER TABLE profile_reports
    ADD CONSTRAINT chk_profile_reports_status
    CHECK (status IN ('open', 'resolved', 'dismissed'));

-- BUG-038: Prevent self-favoriting at database level.
ALTER TABLE user_favorites
    ADD CONSTRAINT chk_no_self_favorite
    CHECK (user_id <> favorite_user_id);
