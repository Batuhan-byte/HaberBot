# Implementasyon Planı: Auth & Role Management

**Tarih:** 2026-05-29  
**Versiyon:** 1.0  
**Durum:** Planlandı

---

## Özet

JWT access token + HttpOnly refresh cookie tabanlı auth, Admin/User rolleri, yorum altyapısı ve seed admin hesabı devreye alınacak. Admin paneli role-based auth ile korunacak, mevcut API key mekanizması ayrı kalacak.

---

## Mevcut Durum Analizi

- Kullanıcı kayıt/giriş ve rol yapısı yok.
- Admin paneli API key ile korunuyor.
- Yorum altyapısı bulunmuyor.

---

## Faz 01: Veri Katmanı (Schema + Config)

### 1.1 Users tablosu
**Hedef:** `users` tablosu (username, password_hash, role, timestamps) ve refresh token alanları.

**Dosyalar:**
- `backend/migrations/007_create_users.up.sql`
- `backend/migrations/007_create_users.down.sql`
- `backend/cmd/migrate/main.go`

---

### 1.2 Comments tablosu
**Hedef:** `comments` tablosu (article_id, user_id, content, created_at).

**Dosyalar:**
- `backend/migrations/008_create_comments.up.sql`
- `backend/migrations/008_create_comments.down.sql`
- `backend/cmd/migrate/main.go`

---

### 1.3 Config alanları
**Hedef:** JWT secret, access/refresh TTL, cookie ayarları ve rate limit parametreleri.

**Dosyalar:**
- `backend/internal/infrastructure/config/config.go`

---

## Faz 02: Domain & Repository

### 2.1 User entity + port
**Hedef:** `User` entity ve `UserRepository` portu.

**Dosyalar:**
- `backend/internal/domain/entity/user.go`
- `backend/internal/domain/port/user_repository.go`

---

### 2.2 Postgres user repo
**Hedef:** Kullanıcı create/find/update (refresh token hash) işlemleri.

**Dosyalar:**
- `backend/internal/adapter/repository/postgres_user.go`
- `backend/internal/adapter/repository/postgres_user_test.go`

---

### 2.3 Comment entity + port + repo
**Hedef:** Comment entity, port ve Postgres repo.

**Dosyalar:**
- `backend/internal/domain/entity/comment.go`
- `backend/internal/domain/port/comment_repository.go`
- `backend/internal/adapter/repository/postgres_comment.go`

---

## Faz 03: Usecase

### 3.1 Auth usecase’leri
**Hedef:** Register, Login, Refresh, Logout usecase’leri; bcrypt hash ve token rotation.

**Dosyalar:**
- `backend/internal/usecase/auth_register.go`
- `backend/internal/usecase/auth_login.go`
- `backend/internal/usecase/auth_refresh.go`
- `backend/internal/usecase/auth_logout.go`

---

### 3.2 Comments usecase’leri
**Hedef:** Comment create/list usecase’leri.

**Dosyalar:**
- `backend/internal/usecase/comments.go`

---

## Faz 04: HTTP Layer & Middleware

### 4.1 Auth handler + middleware
**Hedef:** `/api/auth/*` endpoint’leri ve JWT doğrulama middleware’i.

**Dosyalar:**
- `backend/internal/adapter/handler/auth_handler.go`
- `backend/internal/infrastructure/server/fiber.go`

---

### 4.2 Admin role gate
**Hedef:** Admin handler’larda role kontrolü; API key mekanizması ayrı kalır.

**Dosyalar:**
- `backend/internal/adapter/handler/admin_handler.go`
- `backend/internal/infrastructure/server/fiber.go`

---

### 4.3 Comments handler
**Hedef:** `POST /api/comments` (auth required) ve `GET /api/comments` (public).

**Dosyalar:**
- `backend/internal/adapter/handler/comment_handler.go`
- `backend/internal/infrastructure/server/fiber.go`

---

## Faz 05: Seed + Güvenlik

### 5.1 Seed admin
**Hedef:** `admin1/admin1` admin hesabı bcrypt hash ile oluşturulur.

**Dosyalar:**
- `backend/migrations/009_seed_admin.up.sql` (veya bootstrap kodu)
- `backend/cmd/migrate/main.go`

---

### 5.2 Rate limit + security headers
**Hedef:** Login/register rate limit ve temel security headers.

**Dosyalar:**
- `backend/internal/infrastructure/server/fiber.go`

---

## Görev Listesi

| # | Görev | Dosyalar | Durum |
|---|-------|----------|-------|
| 1 | Users migration | migrations/007_* | Beklemede |
| 2 | Comments migration | migrations/008_* | Beklemede |
| 3 | Config alanları | config.go | Beklemede |
| 4 | User domain + repo | user.go, user_repository.go, postgres_user.go | Beklemede |
| 5 | Comment domain + repo | comment.go, comment_repository.go, postgres_comment.go | Beklemede |
| 6 | Auth usecase’leri | auth_*.go | Beklemede |
| 7 | Auth handler + middleware | auth_handler.go, fiber.go | Beklemede |
| 8 | Admin role gate | admin_handler.go | Beklemede |
| 9 | Comments handler | comment_handler.go | Beklemede |
| 10 | Seed admin | migration/seed | Beklemede |
| 11 | Rate limit + headers | fiber.go | Beklemede |

---

## Success Criteria

- [ ] Auth API çalışır; access token döner, refresh cookie set edilir.
- [ ] User/Admin rol ayrımı ve admin gate çalışır.
- [ ] Şifreler bcrypt hash’li saklanır.
- [ ] Yorum oluşturma ve listeleme endpoint’leri çalışır.
- [ ] Seed admin hesabı oluşur.

---

## Bağımlılıklar

1 → 2 → 3 → 4 → 6 → 7 → 8  
2 → 5 → 9  
1 → 10  
7 → 11

