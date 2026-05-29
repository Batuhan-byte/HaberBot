# Implementation Plan: User Profile (Phase 8)

**Spec:** `docs/superpowers/specs/2026-05-29-user-profile-design.md`  
**Phase:** 8  
**Date Created:** 2026-05-29

---

## Overview

User Profile phase adds public user profile pages, profile editing (bio & avatar), favorites system, and profile reporting. Implementation follows vertical slicing: Foundation → Core API → Frontend → Avatar Storage & Polish.

**Key Dependencies:**
- Auth & Role Management (Phase 7) — must complete first
- Database schema changes (3 new tables)
- File upload handling (avatar storage)
- API + Frontend integration

---

## Architecture Decisions

1. **Avatar Storage:** Local disk (`public/avatars/user-{userId}.[ext]`) for simplicity; CDN migration possible later
2. **Preset Avatars:** Stored in frontend public assets; no database tracking needed
3. **Favorites:** Many-to-many relationship via `user_favorites` table
4. **Reports:** Separate `profile_reports` table; no auto-action (manual admin review)
5. **Edit Permissions:** Only profile owner can edit; enforced at API + frontend level
6. **Profile URL:** `/profile/:username` (URL-friendly, not numeric IDs)
7. **Deleted Accounts:** Profiles archived (remain visible, marked deleted)
8. **Activity Count:** Cached in `user_stats` table, updated on article publish/delete

---

## Task List

### Phase 1: Foundation (Database & Models)

#### Task 1: Database Schema — Core Tables
**Description:** Create database migrations for `user_stats` and `user_favorites` tables. Alter `users` table to add `bio` and `avatar_url` columns. Add indexes for performance.

**Acceptance Criteria:**
- [ ] Users table has `bio` (VARCHAR 160) and `avatar_url` (VARCHAR 255) columns
- [ ] `user_stats` table created with user_id (FK), articles_published, updated_at
- [ ] `user_favorites` table created with (user_id, favorite_user_id) unique constraint
- [ ] All foreign keys configured correctly
- [ ] Migrations are reversible (down scripts work)

**Verification:**
- [ ] `npm run migrate up` succeeds without errors
- [ ] Schema inspection confirms all columns/tables exist
- [ ] Indexes created on (user_id), (user_id, favorite_user_id)
- [ ] Manual DB check: `SELECT * FROM information_schema.TABLES WHERE TABLE_SCHEMA='haberbot'`

**Dependencies:** None (independent)

**Files Likely Touched:**
- `backend/migrations/[timestamp]_create_user_stats.sql`
- `backend/migrations/[timestamp]_alter_users_add_bio_avatar.sql`
- `backend/migrations/[timestamp]_create_user_favorites.sql`

**Estimated Scope:** Small (3 SQL files)

---

#### Task 2: Database Schema — Profile Reports Table
**Description:** Create `profile_reports` table for user moderation. Tracks reports with reason, status, and timestamps.

**Acceptance Criteria:**
- [ ] `profile_reports` table created with id (PK), reporter_user_id (FK), reported_user_id (FK), reason (ENUM), comment, status (ENUM), created_at, resolved_at
- [ ] Reason ENUM: spam, offensive_content, harassment, other
- [ ] Status ENUM: open, reviewed, dismissed, resolved
- [ ] Foreign keys enforce referential integrity
- [ ] Migration is reversible

**Verification:**
- [ ] Migration runs successfully
- [ ] Schema inspection shows all columns/enums
- [ ] Test insert: verify ENUM constraints work

**Dependencies:** Task 1

**Files Likely Touched:**
- `backend/migrations/[timestamp]_create_profile_reports.sql`

**Estimated Scope:** XS (1 SQL file)

---

#### Task 3: TypeScript Models & Interfaces
**Description:** Define TypeScript interfaces for all profile-related data: User profile (public view), ProfileUpdate, UserStats, Favorite, ProfileReport. Generate from database schema.

**Acceptance Criteria:**
- [ ] `User` interface includes id, username, bio, avatar_url, join_date, is_deleted
- [ ] `UserStats` interface includes user_id, articles_published, updated_at
- [ ] `Favorite` interface includes id, user_id, favorite_user_id, created_at
- [ ] `ProfileReport` interface includes id, reporter_user_id, reported_user_id, reason, comment, status, created_at, resolved_at
- [ ] `ProfileUpdateInput` type restricts to bio + avatar_url only
- [ ] All types exported from `src/types/profile.ts`

**Verification:**
- [ ] TypeScript compilation succeeds: `npm run build`
- [ ] No type errors in IDE
- [ ] Interfaces match database columns

**Dependencies:** Task 2

**Files Likely Touched:**
- `backend/src/types/profile.ts`

**Estimated Scope:** XS (1 TS file)

---

#### Task 4: Validation & Sanitization Utilities
**Description:** Create validation functions for profile input: bio length, avatar file type/size, URL slug, etc. Create sanitization for XSS prevention.

**Acceptance Criteria:**
- [ ] `validateBio(bio: string): { valid: boolean; error?: string }`— checks 160 char max, rejects HTML/scripts
- [ ] `validateAvatarFile(file: File | Buffer): { valid: boolean; error?: string }` — checks MIME type (JPG/PNG/WebP), size <2MB
- [ ] `sanitizeBio(bio: string): string` — removes HTML entities, prevents XSS
- [ ] `validateUsername(username: string): { valid: boolean }` — ensures URL-safe slug
- [ ] All functions export from `src/utils/profile-validation.ts`
- [ ] Unit tests for each function

**Verification:**
- [ ] All unit tests pass: `npm test -- profile-validation`
- [ ] Edge cases tested: empty string, 160 char, 161 char, HTML tags, quotes, etc.
- [ ] Manual test: invalid bio rejected, valid bio accepted

**Dependencies:** Task 3

**Files Likely Touched:**
- `backend/src/utils/profile-validation.ts`
- `backend/tests/profile-validation.test.ts`

**Estimated Scope:** Small (2 files)

---

### Checkpoint: Foundation ✓
- [ ] Database migrations run without error
- [ ] TypeScript compilation succeeds
- [ ] All validation unit tests pass
- [ ] Schema matches spec requirements

---

### Phase 2: Backend API Endpoints

#### Task 5: GET /api/users/:username — Profile Fetch
**Description:** Fetch public profile data for a given username. Returns user info + activity stats. Returns 404 if not found, 401 if not authenticated.

**Acceptance Criteria:**
- [ ] Endpoint returns: { id, username, bio, avatar_url, join_date, articles_published, is_deleted }
- [ ] User not authenticated → 401 Unauthorized (redirect to login)
- [ ] Username not found → 404 Not Found
- [ ] Deleted accounts return profile with is_deleted=true
- [ ] Response includes X-Total-Articles header (for caching hints)
- [ ] Endpoint tested with 5+ cases (found, not found, deleted, no auth, etc.)

**Verification:**
- [ ] Integration tests pass: `npm test -- api/users/profile`
- [ ] Manual test: `curl http://localhost:3000/api/users/john-doe` (with auth token)
- [ ] Verify 401 without token, 404 for fake username

**Dependencies:** Task 4

**Files Likely Touched:**
- `backend/src/routes/users.ts` (new endpoint handler)
- `backend/src/services/user-profile.service.ts` (business logic)
- `backend/tests/api/users.profile.test.ts`

**Estimated Scope:** Medium (3 files)

---

#### Task 6: PUT /api/users/profile — Update Own Profile
**Description:** Update authenticated user's bio and/or avatar. Only owner can update. Handles file upload (avatar) and bio text.

**Acceptance Criteria:**
- [ ] Accepts multipart/form-data: bio (text) + avatar (file, optional)
- [ ] Validates bio length (160 max), file size (2MB max), MIME type
- [ ] Only authenticated user can update their own profile (403 if not owner)
- [ ] Updates `users.bio` and `users.avatar_url` atomically
- [ ] Old avatar file deleted from disk if replaced
- [ ] Returns updated user object
- [ ] 400 error if bio >160 chars or avatar >2MB
- [ ] Tested with 8+ cases (valid update, over limit, invalid file, wrong user, etc.)

**Verification:**
- [ ] Integration tests pass: `npm test -- api/users/profile-update`
- [ ] Manual test: upload valid image, verify file on disk
- [ ] Verify 403 when updating another user's profile
- [ ] Verify old file deleted when avatar replaced

**Dependencies:** Task 5

**Files Likely Touched:**
- `backend/src/routes/users.ts` (endpoint)
- `backend/src/services/user-profile.service.ts` (update logic)
- `backend/src/middleware/file-upload.ts` (multer config, if new)
- `backend/tests/api/users.profile-update.test.ts`

**Estimated Scope:** Medium (4 files)

---

#### Task 7: POST /api/users/:userId/favorites — Add to Favorites
**Description:** Add another user to authenticated user's favorites. Returns success/conflict (already favorited).

**Acceptance Criteria:**
- [ ] Requires authentication
- [ ] User cannot favorite themselves (400 Bad Request)
- [ ] User cannot favorite same person twice (409 Conflict: "Already favorited")
- [ ] Creates row in `user_favorites` table
- [ ] Returns { message: "Added to favorites", favorite_user_id: X }
- [ ] Tests: valid add, self-favorite (rejected), duplicate (rejected), no auth (401)

**Verification:**
- [ ] Unit tests pass: `npm test -- api/favorites-add`
- [ ] Manual test: add user to favorites, verify DB row created
- [ ] Verify 409 on duplicate attempt

**Dependencies:** Task 5

**Files Likely Touched:**
- `backend/src/routes/favorites.ts` (new file)
- `backend/src/services/favorites.service.ts` (new file)
- `backend/tests/api/favorites-add.test.ts`

**Estimated Scope:** Small (3 files)

---

#### Task 8: DELETE /api/users/:userId/favorites — Remove from Favorites
**Description:** Remove user from authenticated user's favorites. Returns 204 No Content on success, 404 if not in favorites.

**Acceptance Criteria:**
- [ ] Requires authentication
- [ ] Deletes row from `user_favorites` table
- [ ] Returns 204 No Content on success
- [ ] Returns 404 if favorite relationship doesn't exist
- [ ] Tests: valid delete, not found, no auth

**Verification:**
- [ ] Unit tests pass: `npm test -- api/favorites-remove`
- [ ] Manual test: add favorite, then remove, verify DB row deleted
- [ ] Verify 404 when removing non-existent favorite

**Dependencies:** Task 7

**Files Likely Touched:**
- `backend/src/routes/favorites.ts` (extend Task 7)
- `backend/src/services/favorites.service.ts` (extend Task 7)
- `backend/tests/api/favorites-remove.test.ts`

**Estimated Scope:** XS (extend existing files)

---

#### Task 9: GET /api/users/me/favorites — List My Favorites
**Description:** Fetch authenticated user's favorites list. Paginated, sorted by most-recently-added.

**Acceptance Criteria:**
- [ ] Requires authentication
- [ ] Returns array of favorite users: { id, username, avatar_url, bio, articles_published }
- [ ] Supports pagination (limit, offset query params, defaults 20/0)
- [ ] Sorted by created_at DESC (most recent first)
- [ ] Returns X-Total-Count header (for UI pagination)
- [ ] Tests: valid list, pagination, no auth

**Verification:**
- [ ] Unit tests pass: `npm test -- api/favorites-list`
- [ ] Manual test: `curl http://localhost:3000/api/users/me/favorites?limit=10` (with auth)
- [ ] Verify pagination headers present

**Dependencies:** Task 7

**Files Likely Touched:**
- `backend/src/routes/favorites.ts` (extend)
- `backend/src/services/favorites.service.ts` (extend)
- `backend/tests/api/favorites-list.test.ts`

**Estimated Scope:** Small (3 files)

---

#### Task 10: POST /api/profile-reports — Report Profile
**Description:** Create a report for inappropriate profile. User can report other profiles (not own).

**Acceptance Criteria:**
- [ ] Requires authentication
- [ ] User cannot report themselves (400 Bad Request)
- [ ] Accepts reason (enum: spam, offensive_content, harassment, other) + optional comment
- [ ] Validates reason is one of allowed values
- [ ] Creates row in `profile_reports` table with status='open'
- [ ] Returns { id, message: "Report submitted" }
- [ ] One report per (reporter, reported, reason) per 24hrs (409 Conflict if duplicate)
- [ ] Tests: valid report, self-report, invalid reason, duplicate, no auth

**Verification:**
- [ ] Unit tests pass: `npm test -- api/profile-reports`
- [ ] Manual test: create report, verify DB row
- [ ] Verify 409 on duplicate within 24hrs
- [ ] Verify 400 on invalid reason

**Dependencies:** Task 5

**Files Likely Touched:**
- `backend/src/routes/profile-reports.ts` (new file)
- `backend/src/services/profile-reports.service.ts` (new file)
- `backend/tests/api/profile-reports.test.ts`

**Estimated Scope:** Small (3 files)

---

### Checkpoint: Backend API ✓
- [ ] All 6 endpoints functional (GET profile, PUT update, POST/DELETE favorite, GET favorites list, POST report)
- [ ] All integration tests pass
- [ ] Manual API testing with curl confirms end-to-end flow
- [ ] Error handling verified (401, 403, 404, 400, 409 cases)

---

### Phase 3: Frontend Components

#### Task 11: ProfilePage Component — View Profile
**Description:** React page component `/profile/:username` that displays user profile with all public info, edit button (if own profile), favorites/report buttons (if other profile).

**Acceptance Criteria:**
- [ ] Fetches profile data via GET `/api/users/:username` on mount
- [ ] Displays: avatar, username, bio, join date, articles_published
- [ ] Shows loading spinner while fetching
- [ ] Shows error message if fetch fails (404, 401, server error)
- [ ] If own profile: shows Edit button (visible, enabled)
- [ ] If other profile: shows Add Favorites + Report Profile buttons
- [ ] Edit button opens ProfileEditModal (Task 12)
- [ ] Favorites button toggles favorite (calls POST/DELETE endpoint)
- [ ] Report button opens ReportModal (Task 13)
- [ ] Handles deleted profiles: shows "User deleted" badge, no Edit button
- [ ] Responsive: works on mobile, tablet, desktop
- [ ] Tests: render profile, favorite toggle, report flow, error states

**Verification:**
- [ ] Component tests pass: `npm test -- ProfilePage`
- [ ] Manual browser test: navigate to `/profile/john-doe`, verify all elements render
- [ ] Verify Edit button only shows on own profile
- [ ] Verify API call on mount, data displayed correctly
- [ ] Verify error message on 404

**Dependencies:** Task 5

**Files Likely Touched:**
- `frontend/src/pages/ProfilePage.tsx` (new file)
- `frontend/src/components/ProfileCard.tsx` (reusable card, optional)
- `frontend/tests/ProfilePage.test.tsx`

**Estimated Scope:** Medium (3 files)

---

#### Task 12: ProfileEditModal Component
**Description:** Modal form for editing own profile. Allows bio text + avatar upload/preset selection.

**Acceptance Criteria:**
- [ ] Modal form with:
  - Bio textarea (max 160 chars, live counter showing "X/160")
  - Avatar upload input (accept JPG/PNG/WebP, max 2MB)
  - OR preset avatar selector (grid/dropdown of 7-8 presets)
  - Save & Cancel buttons
- [ ] Validates on client: bio length, file size, MIME type
- [ ] Shows error toast if validation fails
- [ ] Shows loading spinner on Save (disables button)
- [ ] Sends multipart/form-data PUT request to `/api/users/profile`
- [ ] On success: show success toast, close modal, call parent to refresh profile
- [ ] On failure: show error toast, keep modal open
- [ ] Tests: form submission, validation, error handling

**Verification:**
- [ ] Component tests pass: `npm test -- ProfileEditModal`
- [ ] Manual test: open modal, edit bio, upload avatar, save
- [ ] Verify char counter updates live
- [ ] Verify error toast on file >2MB
- [ ] Verify modal closes on success

**Dependencies:** Task 6, Task 12 (profile component)

**Files Likely Touched:**
- `frontend/src/components/ProfileEditModal.tsx` (new file)
- `frontend/tests/ProfileEditModal.test.tsx`

**Estimated Scope:** Medium (2 files)

---

#### Task 13: ReportModal Component & Favorites Button
**Description:** Modal for reporting profile (reason dropdown + comment) + reusable AddToFavorites button.

**Acceptance Criteria:**

**ReportModal:**
- [ ] Modal with dropdown for reason (spam, offensive_content, harassment, other)
- [ ] Optional textarea for comment
- [ ] Submit & Cancel buttons
- [ ] Validates reason is not empty
- [ ] Sends POST `/api/profile-reports` on submit
- [ ] Shows loading spinner on submit
- [ ] On success: show success toast, close modal
- [ ] On error (409 duplicate): show specific message "Already reported"

**AddToFavoritesButton:**
- [ ] Star icon + "Add to Favorites" text
- [ ] Click toggles: POST if not favorited, DELETE if already favorited
- [ ] Button updates immediately (optimistic UI)
- [ ] Shows loading state while request in-flight
- [ ] Shows error toast on failure (403, etc.)
- [ ] Disabled on own profile

**Tests:** modal submission, favorite toggle, error states

**Verification:**
- [ ] Component tests pass: `npm test -- ReportModal ProfileFavoritesButton`
- [ ] Manual test: click Add to Favorites, verify star fills/unfills
- [ ] Manual test: open report modal, submit, verify success toast
- [ ] Verify 409 error handled gracefully

**Dependencies:** Task 10, Task 7–8

**Files Likely Touched:**
- `frontend/src/components/ReportModal.tsx` (new file)
- `frontend/src/components/AddToFavoritesButton.tsx` (new file)
- `frontend/tests/ReportModal.test.tsx`
- `frontend/tests/AddToFavoritesButton.test.tsx`

**Estimated Scope:** Medium (4 files)

---

### Checkpoint: Frontend Core ✓
- [ ] ProfilePage renders correctly
- [ ] ProfileEditModal opens/closes, form validation works
- [ ] AddToFavoritesButton toggles, ReportModal submits
- [ ] All component tests pass
- [ ] Manual end-to-end: view profile → edit → favorite → report

---

### Phase 4: Avatar Storage & Polish

#### Task 14: Avatar File Upload Handler & Storage
**Description:** Implement file upload middleware (multer), validate files, store on disk, generate avatar URL. Handle replace logic (delete old file).

**Acceptance Criteria:**
- [ ] Multer configured for avatar uploads (max 2MB, jpg/png/webp only)
- [ ] File stored at `public/avatars/user-{userId}.{ext}`
- [ ] Old avatar file deleted when replaced (cleanup)
- [ ] Returns file path/URL to client
- [ ] Validates MIME type on server (defense in depth)
- [ ] Non-image files rejected with 400 error
- [ ] Disk full errors handled gracefully (500 error)
- [ ] Tests: valid upload, file too large, invalid type, replace & cleanup

**Verification:**
- [ ] Unit tests pass: `npm test -- avatar-upload`
- [ ] Manual test: upload image via API, verify file on disk
- [ ] Manual test: replace image, verify old file deleted
- [ ] Manual test: upload >2MB file, verify 400 error

**Dependencies:** Task 6

**Files Likely Touched:**
- `backend/src/middleware/multer-avatar.ts` (new file, or extend file-upload)
- `backend/src/services/avatar-storage.service.ts` (new file)
- `backend/tests/avatar-upload.test.ts`

**Estimated Scope:** Small (3 files)

---

#### Task 15: Preset Avatar Setup & Frontend Integration
**Description:** Create/obtain 7–8 preset avatars and place in `public/avatars/presets/`. Expose via simple endpoint. Integrate into ProfileEditModal.

**Acceptance Criteria:**
- [ ] 7–8 PNG preset avatars created/obtained (consistent style, 200x200px, <50KB each)
- [ ] Presets stored in `public/avatars/presets/{name}.png`
- [ ] GET `/api/avatars/presets` returns list of preset URLs: `[ { id, name, url }, ... ]`
- [ ] ProfileEditModal preset selector shows thumbnails
- [ ] User can click preset to select it; selected preset highlighted
- [ ] Save sends preset name (or URL) to backend in avatar_url field
- [ ] Backend recognizes preset URLs vs. uploaded files

**Verification:**
- [ ] Manual test: visit preset endpoint, verify JSON response
- [ ] Manual test: open ProfileEditModal, see preset options, select one, save
- [ ] Verify avatar_url stored correctly in DB (preset URL or upload path)
- [ ] Check preset files exist and serve correctly (200 status)

**Dependencies:** Task 12

**Files Likely Touched:**
- `public/avatars/presets/*.png` (7–8 new image files)
- `backend/src/routes/avatars.ts` (new endpoint)
- `frontend/src/components/ProfileEditModal.tsx` (extend Task 12)

**Estimated Scope:** Small (preset images + 1 endpoint)

---

#### Task 16: Activity Count Cache & Updates
**Description:** Populate `user_stats.articles_published` on user creation. Update count when articles published/deleted. Add cron job or event hook to keep cache fresh.

**Acceptance Criteria:**
- [ ] On user creation: insert row into `user_stats` with articles_published=0
- [ ] When article published (is_approved=true): increment user_stats.articles_published
- [ ] When article deleted: decrement user_stats.articles_published
- [ ] On profile fetch: join with user_stats to get count
- [ ] Migrations populate existing users' counts (backfill)
- [ ] Tests: create user, publish articles, verify count updates

**Verification:**
- [ ] Unit tests pass: `npm test -- user-stats-cache`
- [ ] Manual test: create article as user, verify count increments on profile
- [ ] Manual test: delete article, verify count decrements
- [ ] Check DB: user_stats rows match actual article count

**Dependencies:** Task 1, Task 5

**Files Likely Touched:**
- `backend/src/services/user-profile.service.ts` (extend)
- `backend/src/services/article.service.ts` (extend)
- `backend/migrations/[timestamp]_backfill_user_stats.sql`
- `backend/tests/user-stats.test.ts`

**Estimated Scope:** Medium (4 files)

---

#### Task 17: Integration Tests — End-to-End Profile Flow
**Description:** E2E tests covering full user profile flow: view profile → edit → favorites → report. Tests auth, permissions, data consistency.

**Acceptance Criteria:**
- [ ] Test: User A views User B's profile, sees correct info
- [ ] Test: User A adds User B to favorites, User A's favorite list updates
- [ ] Test: User A edits own profile (bio + avatar), changes persist
- [ ] Test: User A reports User B's profile, report created with correct status
- [ ] Test: User A cannot edit User B's profile (403)
- [ ] Test: User A cannot favorite themselves
- [ ] Test: Deleted user's profile still visible, no Edit button
- [ ] Test: Non-auth user redirected to login
- [ ] All tests mock or use real DB (test database)

**Verification:**
- [ ] All e2e tests pass: `npm test -- profile-e2e`
- [ ] Coverage includes happy path + error cases
- [ ] Manual verification: repeat flows in browser

**Dependencies:** All Tasks 5–16

**Files Likely Touched:**
- `backend/tests/e2e/profile.e2e.test.ts` (new file)

**Estimated Scope:** Medium (1 comprehensive test file)

---

#### Task 18: Documentation & Handoff
**Description:** Update API docs (OpenAPI/Swagger), frontend component docs (Storybook optional), README sections for profile feature. Commit spec + code.

**Acceptance Criteria:**
- [ ] API endpoints documented in OpenAPI spec (or README)
- [ ] Response schemas documented with examples
- [ ] Error codes documented (401, 403, 404, 400, 409)
- [ ] README section added: "User Profiles" with overview + usage examples
- [ ] Component interfaces documented (JSDoc comments)
- [ ] Database schema documented in README/ARCHITECTURE.md
- [ ] Git commit with all changes, spec, and tests

**Verification:**
- [ ] Documentation readable and complete (no placeholders)
- [ ] Examples runnable (curl commands, API calls work)
- [ ] Links/refs correct (no broken paths)

**Dependencies:** All Tasks 1–17

**Files Likely Touched:**
- `docs/API.md` (or OpenAPI spec)
- `README.md`
- `docs/ARCHITECTURE.md` (optional)
- JSDoc comments in components
- Git commit message

**Estimated Scope:** Small (documentation files)

---

### Checkpoint: Complete ✓
- [ ] All tasks complete and tested
- [ ] All endpoints tested with integration tests
- [ ] All components tested with unit/component tests
- [ ] E2E flow validated
- [ ] Documentation complete
- [ ] Code committed to main branch

---

## Task Dependency Graph

```
Foundation Phase:
├── Task 1: User Stats & Favorites Schema
├── Task 2: Profile Reports Schema → Task 1
├── Task 3: TypeScript Models → Task 2
└── Task 4: Validation Utils → Task 3

Backend API Phase:
├── Task 5: GET Profile → Task 4
├── Task 6: PUT Update Profile → Task 5
├── Task 7: POST Add Favorite → Task 5
├── Task 8: DELETE Remove Favorite → Task 7
├── Task 9: GET Favorites List → Task 7
└── Task 10: POST Report Profile → Task 5

Frontend Phase:
├── Task 11: ProfilePage Component → Task 5
├── Task 12: ProfileEditModal Component → Task 6, Task 11
├── Task 13: ReportModal & Favorites Button → Task 10, Task 7–8
├── Task 14: Avatar Upload Handler → Task 6
└── Task 15: Preset Avatar Setup → Task 12

Polish & Integration Phase:
├── Task 16: Activity Count Cache → Task 1, Task 5
├── Task 17: E2E Tests → All Tasks 5–16
└── Task 18: Documentation & Handoff → All Tasks
```

---

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Avatar file collisions or disk space exhaustion | Medium | Use UUID for filenames, implement cleanup cron job, monitor disk usage |
| XSS via bio field or report comments | High | Sanitize all user input server-side, use parameterized queries, CSP headers |
| Profile data consistency issues (deleted user, stale cache) | Medium | Use DB transactions for updates, periodic cache refresh, soft deletes with archiving |
| Slow profile page load (N+1 queries) | Medium | Join user_stats, favorites in single query, cache results, add indexes |
| Avatar upload abuse (spam, large files) | Medium | Rate limit uploads, validate MIME type, store outside web root initially |
| Concurrent favorite/report submissions (race condition) | Low | DB unique constraints, idempotent endpoints, test concurrency scenarios |

---

## Open Questions

1. **Preset Avatar Source:** Should we create custom avatars or use an existing avatar library (Avataaars, Dicebear, etc.)?
2. **Avatar Compression:** Should we auto-compress uploaded avatars to save disk space?
3. **Follow-Up Features:** Should "My Favorites" have a dedicated page, or is list view enough for now?
4. **Profile Analytics:** Should we track profile view counts (future phase)?
5. **Username Changes:** Can users change their username, or is it immutable?

---

## Success Criteria (Phase Complete)

✅ All 18 tasks complete  
✅ All unit + integration + E2E tests passing  
✅ API endpoints functional and documented  
✅ Frontend components responsive and tested  
✅ Avatar storage working (upload, preset, delete)  
✅ Favorites and reporting systems operational  
✅ Database migrations applied successfully  
✅ Code merged to main branch  
✅ Documentation up-to-date  
✅ Zero ambiguity (spec score 0.10)

---

**Next Phase:** Phase 9 (TBD) — could be comments integration, advanced search, or user discovery features.
