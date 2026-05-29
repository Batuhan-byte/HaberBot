# User Profile Design Specification
**Date:** 2026-05-29  
**Status:** Approved  
**Ambiguity Score:** 0.10

---

## 1. Overview

The User Profile feature allows each user on HaberBot to have a public-facing profile page displaying their information, activity history, and social connections. Profiles are viewable by all authenticated news publishers on the platform and are editable only by the profile owner.

---

## 2. Core Features

### 2.1 Profile Information

Each profile displays the following information:

| Field | Type | Editable | Visible | Max Length |
|-------|------|----------|---------|-----------|
| Username | String | No | Public | — |
| Bio | Text | Yes (owner only) | Public | 160 characters |
| Avatar | Image | Yes (owner only) | Public | 2 MB |
| Join Date | Date | No | Public | — |
| Activity Count | Integer | No | Public | — |

**Activity Count:** Number of news articles published by the user (read-only, computed from news table).

---

### 2.2 Avatar Management

#### Upload
- **Accepted Formats:** JPG, PNG, WebP
- **Max Size:** 2 MB
- **Storage:** Local disk at `public/avatars/user-{userId}.[ext]`
- **Compression:** Optional (recommended for perf)
- **Validation:** File type check on upload; reject files >2MB or unsupported formats

#### Preset Avatars
- **Count:** 7–8 default avatars provided by the platform
- **Location:** `public/avatars/presets/{preset-name}.png`
- **Fallback:** If user has no uploaded avatar, show first preset by default
- **User Selection:** Dropdown or grid in profile edit form

#### Replace Logic
- User can upload custom avatar to replace current avatar (uploaded or preset)
- Old uploaded avatar file deleted from disk when replaced
- User can revert to any preset avatar at any time

---

### 2.3 Bio Management

- **Max Length:** 160 characters (enforced client-side and server-side)
- **Format:** Plain text (no HTML, markdown, or formatting)
- **Editing:** Available in profile edit form; visible immediately after save
- **Validation:** No URL/link detection; sanitize to prevent XSS

---

### 2.4 Edit Permissions & UI

- **Who Can Edit:** Only the profile owner
- **Edit Button Visibility:** Shows only when viewing own profile (authenticated user === profile owner)
- **Edit Form:** Modal or separate page with fields:
  - Bio textarea (160 char counter)
  - Avatar upload + preset selector
  - Save & Cancel buttons
- **Success Feedback:** Toast notification on successful save; redirect to updated profile

---

### 2.5 Access Control

#### Public Profiles
- **URL:** `/profile/:username`
- **Viewable By:** All authenticated users (logged-in)
- **Not Logged In:** 401 Unauthorized → redirect to login
- **Non-existent Profile:** 404 Not Found

#### Deleted Accounts
- Profile remains archived (visible) but marked as deleted
- Username remains unique (cannot be reused)
- Avatar and bio remain displayed
- Activity count frozen at deletion time
- All links to profile still work (no 404)
- Optional: Add "User deleted" badge or disabled state

---

### 2.6 Activity Metrics

**Displayed on Profile:**
- **Articles Published:** Count of all news articles published by this user
- **Join Date:** ISO 8601 date (e.g., "May 29, 2026")
- **Last Active:** Optional timestamp of last article published (or last login)

**Calculation:**
- Query `news` table for `author_id = user_id` and `status = "published"`
- Cache count in `user_stats` table, updated whenever article is published/deleted

---

### 2.7 Favorites System

Users can add other users to their favorites list.

**Database:**
- New table: `user_favorites` (user_id, favorite_user_id, created_at)
- Unique constraint: (user_id, favorite_user_id) — no duplicate favorites

**UI:**
- "Add to Favorites" button on other users' profiles
- Button toggles to "Remove from Favorites" if already favorited
- Star icon + text label
- Toast feedback on add/remove
- No dedicated "My Favorites" page in this phase (future scope)

---

### 2.8 Profile Reporting

Users can report inappropriate profiles.

**Database:**
- New table: `profile_reports` (id, reporter_user_id, reported_user_id, reason, status, created_at, resolved_at)
- Reason field: Enum (spam, offensive_content, harassment, etc.)

**UI:**
- "Report Profile" button or link on other users' profiles
- Opens modal with reason dropdown
- Optional comment field
- Submit button; success toast
- User cannot report own profile

**Admin Handling:**
- Admin panel lists reports (separate from this feature, not in this phase)
- Reports marked as resolved/dismissed
- No automatic profile deletion or suspension (manual admin review)

---

## 3. Data Model

### Users Table (Existing)
```sql
ALTER TABLE users ADD COLUMN bio VARCHAR(160) DEFAULT NULL;
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255) DEFAULT NULL;
ALTER TABLE users ADD COLUMN join_date DATETIME DEFAULT CURRENT_TIMESTAMP;
```

### New Tables

**user_stats** (for activity caching)
```sql
CREATE TABLE user_stats (
  user_id INT PRIMARY KEY,
  articles_published INT DEFAULT 0,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
```

**user_favorites**
```sql
CREATE TABLE user_favorites (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  favorite_user_id INT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY (user_id, favorite_user_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (favorite_user_id) REFERENCES users(id)
);
```

**profile_reports**
```sql
CREATE TABLE profile_reports (
  id INT PRIMARY KEY AUTO_INCREMENT,
  reporter_user_id INT NOT NULL,
  reported_user_id INT NOT NULL,
  reason ENUM('spam', 'offensive_content', 'harassment', 'other') DEFAULT 'other',
  comment VARCHAR(500),
  status ENUM('open', 'reviewed', 'dismissed', 'resolved') DEFAULT 'open',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  resolved_at DATETIME,
  FOREIGN KEY (reporter_user_id) REFERENCES users(id),
  FOREIGN KEY (reported_user_id) REFERENCES users(id)
);
```

---

## 4. Architecture

### Frontend Components

1. **ProfilePage** (`src/pages/ProfilePage.tsx`)
   - Fetches user data from `/api/users/:username`
   - Displays profile info, stats, avatar
   - Conditionally shows Edit/Add Favorites/Report buttons
   - Handles favorite toggle & report modal

2. **ProfileEditModal** (`src/components/ProfileEditModal.tsx`)
   - Form for bio (textarea + char counter)
   - Avatar upload (drop zone or file input)
   - Preset avatar selector (grid or dropdown)
   - Submit/Cancel buttons
   - Validation & error handling

3. **ProfileCard** (reusable, optional)
   - Compact profile display for lists
   - Used elsewhere on platform (future)

### Backend Endpoints

| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/api/users/:username` | Optional* | Fetch profile data |
| PUT | `/api/users/profile` | Required | Update own profile (bio, avatar) |
| POST | `/api/users/:userId/favorites` | Required | Add user to favorites |
| DELETE | `/api/users/:userId/favorites` | Required | Remove from favorites |
| GET | `/api/users/:userId/favorites` | Required | Get my favorites list |
| POST | `/api/profile-reports` | Required | Report profile |

*Optional auth: return 401 if not logged in instead of 404. Let user see username but redirect on data fetch.

### File Storage

```
public/
  avatars/
    user-{userId}.jpg
    user-{userId}.png
    user-{userId}.webp
    presets/
      avatar-1.png
      avatar-2.png
      ... (7-8 presets)
```

---

## 5. User Flows

### View Profile
1. User clicks on another user's name/profile link
2. Browser navigates to `/profile/:username`
3. Page fetches profile data via GET `/api/users/:username`
4. If not logged in → 401 → redirect to login, return after auth
5. Display profile: avatar, bio, join date, article count
6. Show Add Favorites & Report buttons (if viewing other user's profile)

### Edit Own Profile
1. On own profile page, click Edit button
2. Modal opens: bio textarea + avatar upload/preset selector
3. User updates bio and/or avatar
4. Click Save
5. PUT request to `/api/users/profile` with bio + avatar
6. Server validates, saves to DB, moves file if upload
7. Success toast; modal closes; page refreshes profile data
8. User sees updated profile immediately

### Add to Favorites
1. On other user's profile, click "Add to Favorites"
2. POST `/api/users/:userId/favorites`
3. Button toggles to "Remove from Favorites"
4. Toast: "User added to favorites"

### Report Profile
1. On other user's profile, click "Report Profile"
2. Modal opens: dropdown (reason), optional comment
3. User selects reason and optionally adds comment
4. Click Submit
5. POST `/api/profile-reports` with reason + comment
6. Server validates, inserts report
7. Success toast; modal closes
8. User cannot report same profile twice in X hours (optional rate limit)

---

## 6. Error Handling

| Scenario | Response |
|----------|----------|
| Profile not found | 404 Not Found |
| Not authenticated | 401 Unauthorized (redirect to login) |
| Avatar upload >2MB | 400 Bad Request: "File too large" |
| Avatar unsupported format | 400 Bad Request: "Unsupported file type" |
| Bio exceeds 160 chars | 400 Bad Request: "Bio must be ≤160 characters" |
| XSS attempt in bio | 400 Bad Request: sanitized; reject if suspicious |
| Favorite already exists | 409 Conflict: "Already in favorites" |
| Report already filed | 409 Conflict: "Already reported this profile" |
| Disk full (avatar upload) | 500 Internal Server Error: "Upload failed" |

---

## 7. Security & Validation

### Input Validation
- Bio: max 160 chars; sanitize HTML entities; reject URLs/scripts
- Avatar: check MIME type + file extension; reject executable files
- Usernames: immutable; validated at account creation

### Authorization
- Edit endpoint: verify `req.user.id === profile_owner_id`
- Favorite/Report: verify `req.user.id` is authenticated and not the reported user

### File Handling
- Store uploads outside `public/` webroot if possible; serve via secure endpoint
- Use UUIDs or hashed userIds for filenames to prevent enumeration
- Implement file size limits at middleware level
- Delete old avatar files when replaced

### Rate Limiting
- Report endpoint: 1 report per profile per user per 24 hours
- Favorite toggle: no limit (user preference)

---

## 8. Testing

### Unit Tests
- Bio validation: length, sanitization
- Avatar validation: MIME type, size
- Favorites toggle: add/remove logic
- Report creation: duplicate prevention

### Integration Tests
- Full edit profile flow (upload avatar + update bio)
- Favorite add/remove with DB state
- Report submission with duplicate check
- Deleted account profile visibility

### Manual Tests
- Avatar upload with various file types (including invalid)
- Bio character counter live update
- Preset avatar selection
- Favorites list accuracy
- Report modal + submission
- Non-existent profile 404
- Unauthenticated access redirects

---

## 9. Future Scope (Out of Phase)

- Profile follow/unfollow (with notifications)
- User activity feed (recent articles, activity timeline)
- Follower/following lists
- Profile customization (themes, layout)
- User badges/achievements
- Profile analytics (view count, engagement metrics)
- Private profiles (visibility toggle)
- Two-factor authentication for profile editing

---

## 10. Summary

User profiles are the foundation of social identity on HaberBot. This phase establishes:
- ✅ Public, read-only profile viewing for authenticated users
- ✅ Owner-editable bio and avatar
- ✅ Activity metrics (article count, join date)
- ✅ Favorites system for discovering users
- ✅ Reporting system for moderation
- ✅ Archived profile preservation on account deletion

**Ambiguity:** Very low (0.10). All UI flows, data models, and permissions clearly defined.
