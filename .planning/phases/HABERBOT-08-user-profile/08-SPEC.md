# User Profile Design Specification (Phase 8)

**Date:** 2026-05-29  
**Phase:** 8  
**Spec File:** `docs/superpowers/specs/2026-05-29-user-profile-design.md`  
**Status:** Approved  
**Ambiguity Score:** 0.10

---

## Quick Summary

Phase 8 adds public user profiles with bio + avatar, profile editing, favorites, and reporting system. Profiles are visible to authenticated users, editable only by owner. Implementation plan: 18 tasks across Foundation → API → Frontend → Polish phases.

---

## Core Requirements (Locked)

1. ✅ Profile page at `/profile/:username` — public, auth-required
2. ✅ Profile shows: username, bio (160 char max), avatar, join date, article count
3. ✅ Profile editing: bio + avatar (upload or preset), owner-only
4. ✅ Avatar upload: JPG/PNG/WebP, max 2MB, stored locally
5. ✅ Preset avatars: 7-8 built-in options in frontend assets
6. ✅ Favorites system: users can add other users to favorites
7. ✅ Profile reporting: users can report inappropriate profiles
8. ✅ Deleted accounts: profiles archived (remain visible, marked deleted)
9. ✅ Activity metrics: article count cached, updated on publish/delete

---

## Architecture & Key Decisions

- **Database:** 3 new tables (user_stats, user_favorites, profile_reports)
- **Avatar Storage:** Local disk (`public/avatars/user-{userId}.[ext]`)
- **Edit Permissions:** Only owner can edit; enforced at API + frontend
- **Favorites:** Many-to-many via `user_favorites` table
- **Reports:** Tracked but no auto-action (manual admin review)
- **Deleted Accounts:** Soft delete — profile archived, not hard-deleted

---

## Files & Documentation

- **Implementation Plan:** `.planning/phases/HABERBOT-08-user-profile/IMPLEMENTATION-PLAN.md` (18 tasks)
- **Full Spec:** `docs/superpowers/specs/2026-05-29-user-profile-design.md`

---

## Success Criteria

- [x] Spec approved by user (all 9 core requirements locked)
- [x] Implementation plan created (18 tasks, clear dependencies)
- [ ] All 18 tasks implemented & tested
- [ ] API endpoints functional (6 endpoints: GET profile, PUT update, POST/DELETE favorite, GET favorites, POST report)
- [ ] Frontend components complete (ProfilePage, ProfileEditModal, ReportModal, FavoritesButton)
- [ ] Avatar upload/storage working
- [ ] Database migrations applied
- [ ] E2E tests passing
- [ ] Documentation complete

---

## Risks & Mitigations

| Risk | Mitigation |
|------|-----------|
| Avatar disk space exhaustion | Cleanup cron, size monitoring |
| XSS via bio/comments | Server-side sanitization, parameterized queries |
| Profile data consistency (cache) | DB transactions, periodic refresh, soft deletes |
| N+1 query performance | JOIN user_stats, add indexes |
| Avatar upload abuse | Rate limiting, MIME validation, outside webroot |

---

## Next Steps

1. Phase 8 initialization: `gsd-tools phase init 08`
2. Assign tasks to implementation agent
3. Checkpoint after Foundation phase (Tasks 1–4)
4. Checkpoint after Backend API phase (Tasks 5–10)
5. Checkpoint after Frontend phase (Tasks 11–13)
6. Final checkpoint after polish (Tasks 14–18)
7. Commit to main branch with all tests passing
