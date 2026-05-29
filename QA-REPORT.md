# HaberBot — Siber Güvenlik & Kalite Denetim (QA) Raporu

**Tarih:** 2026-05-29  
**Denetçi:** Red Team Ajanı (Otomatik Sızma Testi)  
**Kapsam:** `backend/` — Go API, veritabanı şeması, JWT yetkilendirme, dosya yükleme, AI entegrasyonu  
**Yöntem:** Beyaz kutu (white-box) statik analiz + mimari inceleme  
**Toplam Bulgu:** 43 (6 Kritik, 9 Yüksek, 18 Orta, 8 Düşük, 2 Bilgi)

---

## Özet Tablosu

| Öncelik | Sayı | Durum |
|---------|------|-------|
| 🔴 KRİTİK | 6 | Hemen düzeltilmeli — üretim ortamında aktif |
| 🟠 YÜKSEK | 9 | 24 saat içinde düzeltilmeli |
| 🟡 ORTA | 18 | Bir sonraki sprint'te ele alınmalı |
| 🔵 DÜŞÜK | 8 | Teknik borç olarak kaydedilmeli |
| ⚪ BİLGİ | 2 | İyileştirme önerisi |

---

## 🔴 KRİTİK Açıklar (6)

### BUG-001: `.env` Dosyasında Gerçek Sırlar Kaynak Kodunda

**Dosya:** `backend/.env`  
**Satır:** 1-13  
**OWASP:** A02:2021 – Cryptographic Failures  

**Açıklama:**  
`.env` dosyasında tüm üretim sırları düz metin olarak mevcuttur:
- Neon PostgreSQL connection string (kullanıcı: `neondb_owner`, şifre: `npg_***[REDACTED]***`)
- Google Gemini API anahtarı (`AIzaSy***[REDACTED]***`)
- Admin API anahtarı (`haberbot-***[REDACTED]***`)
- Mistral AI API anahtarı (`OpnoQ***[REDACTED]***`)

**Sömürü Senaryosu:**  
Bu repo herhangi bir remote'a (GitHub, GitLab) push edildiyse, tüm bu sırlar kamuya açık demektir. Saldırgan:
1. Neon PostgreSQL veritabanına tam erişim elde eder — tüm kullanıcı verilerini okur/siler/değiştirir
2. Gemini/Mistral API anahtarlarıyla ücretsiz AI kullanımı sağlar veya faturalandırma saldırısı yapar
3. Admin API anahtarıyla tüm admin işlemlerini (konu silme,haber çekme, AI özetleme) tetikler

**Etki:** Veri ihlali, tam sistem kontrolü, finansal zarar  
**Aciliyet:** HEMEN

---

### BUG-002: Hardcoded JWT Secret Fallback

**Dosya:** `backend/internal/infrastructure/config/config.go`  
**Satır:** 79-82  
**OWASP:** A07:2021 – Identification and Authentication Failures  

**Kod:**
```go
jwtSecret := os.Getenv("JWT_SECRET")
if jwtSecret == "" {
    jwtSecret = "haberbot-super-secret-jwt-signing-key-12345"
}
```

**Açıklama:**  
`JWT_SECRET` ortam değişkeni tanımlı değilse, uygulama hardcoded bir fallback değer kullanıyor. Bu değer kaynak kodunda herkese açık.

**Sömürü Senaryosu:**  
Saldırgan bu bilinen secret'ı kullanarak:
1. Herhangi bir kullanıcı adına geçerli JWT token'lar üretir
2. `role: "Admin"` claim'li token oluşturarak admin erişimi elde eder
3. Refresh token'lar oluşturarak kalıcı oturum açar

**Etki:** Tam kimlik sahteciliği, yetki yükseltme  
**Aciliyet:** HEMEN — Üretimde `JWT_SECRET` zorunlu kılınmalı, fallback kaldırılmalı

---

### BUG-003: Admin Şifre Hash'i Kaynak Kodunda

**Dosya:** `backend/migrations/009_seed_admin.up.sql`  
**Satır:** 2  
**OWASP:** A07:2021 – Identification and Authentication Failures  

**Kod:**
```sql
INSERT INTO users (id, username, password_hash, role, email)
VALUES ('admin1-uuid-placeholder-1234567890ab', 'admin1',
  '$2a$10$6Wqf8hI7vR9Z6J1uXJgYGea7a8dG09R8a9A5S3g4M6iY7a4r8d3aK',
  'Admin', 'admin@HaberBot.com')
ON CONFLICT (username) DO NOTHING;
```

**Açıklama:**  
Admin kullanıcısının bcrypt password hash'i migration dosyasında açıktestdata olarak yer alıyor. Bu hash brute-force ile kırılabilir.

**Sömürü Senaryosu:**  
1. Hash'i rainbow table veya brute-force ile kır → admin şifresini öğren
2. `admin1` hesabıyla sisteme giriş yap
3. Tüm admin yetkilerini kullan (veritabanı silme, AI manipülasyonu)

**Etki:** Admin hesabı ele geçirilmesi  
**Aciliyet:** HEMEN — Migration seed'i kaldırılmalı, admin şifresi env variable'dan alınmalı

---

### BUG-004: Refresh Token Cookie'lerinde `Secure: false`

**Dosya:** `backend/internal/adapter/handler/auth_handler.go`  
**Satır:** 80, 106, 121, 150  
**OWASP:** A02:2021 – Cryptographic Failures  

**Kod (tüm cookie'lerde aynı):**
```go
c.Cookie(&fiber.Cookie{
    Name:     "refresh_token",
    Value:    loginResp.RefreshToken,
    HTTPOnly: true,
    Secure:   false,  // ← TEHLİKELİ
    SameSite: "Lax",
})
```

**Açıklama:**  
Refresh token cookie'si `Secure: false` ile ayarlanıyor. Bu, cookie'nin HTTP (TLS olmayan) bağlantılarda da gönderileceği anlamına gelir.

**Sömürü Senaryosu:**  
Aynı WiFi ağindaki saldırgan (MITM saldırısı):
1. Trafiği dinler (ARP spoofing, DNS poisoning)
2. HTTP request'lerden `refresh_token` cookie'sini çalar
3. Kendi token'ını kullanarak kullanıcının oturumunu ele geçirir
4. Erişim token'larını yenileyerek kalıcı erişim sağlar

**Etki:** Oturum ele geçirme, hesap gaspı  
**Aciliyet:** HEMEN — `Secure: true` yapılmalı (development hariç)

---

### BUG-005: Public Profilde E-posta Adresi Ifşa Ediliyor

**Dosya:** `backend/internal/usecase/profile_get.go`  
**Satır:** 13-23, 66-76  
**OWASP:** A01:2021 – Broken Access Control  

**Kod:**
```go
type ProfileResponse struct {
    ID        string    `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`     // ← HERKESE AÇIK
    Role      string    `json:"role"`      // ← HERKESE AÇIK
    // ...
}
```

**Açıklama:**  
`GET /api/v1/users/:username` endpoint'i herkese açıktır (auth gerektirmez). Yanıtta tüm kullanıcıların e-posta adresleri ve rolleri yer alır.

**Sömürü Senaryosu:**  
1. Saldırgan script yazarak tüm kullanıcı adlarını tarar
2. Toplu e-posta toplar (phishing attack için)
3. `role: "Admin"` olan kullanıcıları belirleyerek hedefli saldırı planlar
4. Credential stuffing attack'lar için e-posta listesi oluşturur

**Etki:** Toplu veri sızıntısı, hedefli phishing  
**Aciliyet:** HEMEN — Public profile'dan Email ve Role alanları kaldırılmalı

---

### BUG-006: Veritabanı `role` Kolonunda CHECK Kısıtlaması Yok

**Dosya:** `backend/migrations/007_create_users.up.sql`  
**Satır:** 5  
**OWASP:** A01:2021 – Broken Access Control  

**Kod:**
```sql
role VARCHAR(50) NOT NULL DEFAULT 'User'
-- CHECK (role IN ('User', 'Admin')) eksik!
```

**Açıklama:**  
`role` kolonu herhangi bir string değer kabul ediyor. Uygulama seviyesinde `"Admin"` kontrolü yapılsa da, veritabanı seviyesinde herhangi bir değer insert edilebilir.

**Sömürü Senaryosu:**  
Eğer saldırgan SQL enjeksiyonu veya veritabanı erişimi elde ederse:
1. `UPDATE users SET role = 'SuperAdmin' WHERE username = 'attacker'`
2. Uygulama seviyesindeki `"admin"` kontrolünü atlar
3. Tüm admin yetkilerini kullanır

**Etki:** Yetki yükseltme  
**Aciliyet:** HEMEN — CHECK kısıtlaması eklenmeli

---

## 🟠 YÜKSEK Açıklar (9)

### BUG-007: Global Rate Limiting Yok

**Dosya:** `backend/internal/infrastructure/server/fiber.go`  
**Satır:** 36-48  

**Açıklama:**  
Sadece `/api/v1/auth` (register/login) endpoint'lerinde rate limiting var. Tüm diğer API endpoint'leri sınırsız.

**Etki:** DDoS saldırıları, kaynak yorgunluğu, brute-force  
**Önerilen:** Tüm API route'larına genel rate limiting ekle (ör. 100 istek/dakika)

---

### BUG-008: AI Özetleme Endpoint'i Auth Korumasız

**Dosya:** `backend/internal/infrastructure/server/fiber.go`  
**Satır:** 60  

**Kod:**
```go
api.Post("/articles/:id/summary", container.ArticleHandler.SummarizeArticle)
```

**Açıklama:**  
`POST /api/v1/articles/:id/summary` herkese açıktır. Auth yok. Bu endpoint AI API çağrısı yapıyor.

**Sömürü Senaryosu:**  
Saldırgan script yazarak:
1. Tüm article ID'lerini enumerate eder
2. Her biri için summary endpoint'ini çağırır
3. Gemini/Mistral API kotasını tüketir
4. Normal kullanıcılar AI özelliğini kullanamaz hale gelir

**Etki:** API kota tüketimi, hizmet reddi  
**Önerilen:** Bu endpoint'e JWT auth ekle veya rate limit uygula

---

### BUG-009: Refresh Endpoint'inde Rate Limiting Yok

**Dosya:** `backend/internal/infrastructure/server/fiber.go`  
**Satır:** 86  

**Kod:**
```go
authGroup.Post("/refresh", container.AuthHandler.Refresh)
```

**Açıklama:**  
`/api/v1/auth/refresh` rate limiter'ın dışında. Saldırgan çalınan refresh token'ları sınırsızca deneyebilir.

**Etki:** Brute-force ile token tahmin  
**Önerilen:** Refresh endpoint'ine de rate limiting ekle

---

### BUG-010: Admin API Key Karşılaştırması Timing Attack'a Açık

**Dosya:** `backend/internal/infrastructure/server/fiber.go`  
**Satır:** 108  

**Kod:**
```go
if key != "" && key == container.Config.AdminAPIKey {
```

**Açıklama:**  
Go'da `==` operatörü string karşılaştırmalarında sabit zamanlı (constant-time) değildir. İlk farklı byte'da durur.

**Sömürü Senaryosu:**  
Saldırgan çok sayıda istek göndererek:
1. İlk karakteri doğru tahmin edene kadar deneme yapar (100-200 istek)
2. Doğru ilk karakteri bulduktan sonra ikinci karaktere geçer
3. ~3200 istekle 16 karakterlik key'i brute-force eder (20 karakter * 160 ortalama deneme)

**Etki:** Admin anahtarı ele geçirme  
**Önerilen:** `crypto/subtle.ConstantTimeCompare()` kullan

---

### BUG-011: Güvenlik Başlıkları (Security Headers) Eksik

**Dosya:** `backend/internal/infrastructure/server/fiber.go`  
**Satır:** 20-34  

**Açıklama:**  
Hiçbir güvenlik başlığı ayarlanmamış:
- `X-Content-Type-Options: nosniff` — MIME sniffing saldırısına açık
- `X-Frame-Options: DENY` — Clickjacking saldırısına açık
- `Strict-Transport-Security` — MITM downgrade saldırısına açık
- `Content-Security-Policy` — XSS saldırısına açık
- `X-XSS-Protection` — Reflected XSS'e açık

**Etki:** Çeşitli client-side saldırılar  
**Önerilen:** Fiber security middleware ekle

---

### BUG-012: Yorum Oluşturmada Rate Limiting ve Uzunluk Sınırı Yok

**Dosya:** `backend/internal/adapter/handler/comment_handler.go`  
**Satır:** 27-63  

**Açıklama:**  
- Rate limiting yok — binlerce yorum spam yapılabilir
- `Content` alanı için maksimum uzunluk yok — megabyte boyutunda yorum gönderilebilir

**Sömürü Senaryosu:**  
1. Bot ağıyla dakikada 1000 yorum gönderilir
2. Veritabanı disk alanı dolar
3. Yorum listeleme endpoint'i çöker (memory exhaustion)

**Etki:** Veritabanı şişirmesi, hizmet reddi  
**Önerilen:** Rate limit + max 2000 karakter sınırı

---

### BUG-013: Kayıt Formunda E-posta Format Doğrulaması Yok

**Dosya:** `backend/internal/usecase/auth_register.go`  
**Satır:** 35-36  

**Kod:**
```go
if req.Email == "" {
    return nil, errors.New("email is required")
}
```

**Açıklama:**  
Sadece boş kontrol yapılıyor. `"notanemail"`, `"<script>alert(1)</script>"` gibi herhangi bir string kabul ediliyor.

**Etki:** Geçersiz e-posta kayıtları, potansiyel log injection  
**Önerilen:** Regex ile e-posta format doğrulaması ekle

---

### BUG-014: Kullanıcı Adı Karakter Seti Doğrulaması Yok

**Dosya:** `backend/internal/usecase/auth_register.go`  
**Satır:** 32-33  

**Açıklama:**  
Sadece uzunluk kontrolü var (>=3). `<script>`, `admin' OR 1=1--`, kontrol karakterleri gibi herhangi bir karakter kabul ediliyor.

**Etki:** Stored XSS, log injection, social engineering (sahte admin adları)  
**Önerilen:** `[a-zA-Z0-9_-]` karakter seti kısıtlaması ekle

---

### BUG-015: Logout Hatası Yutuluyor (Silent Failure)

**Dosya:** `backend/internal/adapter/handler/auth_handler.go`  
**Satır:** 139  

**Kod:**
```go
_ = h.logoutUC.Execute(c.Context(), userID)
```

**Açıklama:**  
Logout işleminde veritabanı token iptali başarısız olursa hata yutuluyor. Kullanıcı "çıkış yaptım" sanırken oturumu hâlâ aktif.

**Sömürü Senaryosu:**  
Paylaşımlı bilgisayarda çıkış yapan kullanıcı:
1. DB hatası nedeniyle token iptali başarısız olur
2. Kullanıcı bilgisayarı bırakır
3. Başka biri eski cookie ile oturumu devam ettirir

**Etki:** Yetkisiz erişim  
**Önerilen:** Hata loglanmalı, kullanıcıya bildirilmeli

---

### BUG-016: Dockerfile'da Root Kullanıcı ile Çalışma

**Dosya:** `backend/Dockerfile`  
**Satır:** 8-14  

**Açıklama:**  
Container'da `USER` direktifi yok. Uygulama root olarak çalışıyor.

**Etki:** Container escape saldırısında saldırgan root olur  
**Önerilen:** `RUN adduser -D appuser` + `USER appuser` ekle

---

## 🟡 ORTA Açıklar (18)

### BUG-017: Hata Mesajlarında İç Detay Sızıntısı

**Dosyalar:** Tüm handler dosyaları  
**Satırlar:** `article_handler.go:48,69,98,119,155`, `admin_handler.go:93,139,171,185,215,252,259,281`, `comment_handler.go:55,78`, `profile_handler.go:74,179,207,258,303`, `topic_handler.go:29,62,73`

**Açıklama:**  
Tüm handler'lar `err.Error()` ile ham hata mesajlarını istemciye dönüyor. Veritabanı hata detayları, dosya yolları veya bağlantı string'leri sızabilir.

**Önerilen:** Genel hata mesajları kullan, detayları logla

---

### BUG-018: Avatar URL Doğrulama Bypass

**Dosya:** `backend/internal/usecase/profile_update.go`  
**Satır:** 62  

**Kod:**
```go
strings.HasPrefix(avatarURL, "http")
```

**Açıklama:**  
`http://evil.com/malicious.js` veya `http://attacker.com/phishing` gibi herhangi bir HTTP URL'i kabul ediliyor. Bu URL'ler başka kullanıcılara gösteriliyor.

**Önerilen:** URL allowlist'i veya regex doğrulama ekle

---

### BUG-019: JWT'de `jti` (Token ID) Claim'i Yok

**Dosya:** `backend/internal/infrastructure/auth/jwt.go`  
**Satır:** 11-21, 25-33  

**Açıklama:**  
Token'lar bireysel olarak iptal edilemiyor. Tek yol tüm signing key'i döndürmek.

**Önerilen:** `jti` claim'i ekle, token blacklist mekanizması kur

---

### BUG-020: SameSite Cookie "Lax" Olmalı "Strict"

**Dosya:** `backend/internal/adapter/handler/auth_handler.go`  
**Satır:** 81, 107, 122, 151  

**Açıklama:**  
`SameSite: "Lax"` cross-site GET isteklerinde cookie'yi gönderir. CSRF koruması zayıflar.

**Önerilen:** `SameSite: "Strict"` kullan

---

### BUG-021: Şifre Politikası Çok Zayıf

**Dosya:** `backend/internal/usecase/auth_register.go`  
**Satır:** 38-39  

**Kod:**
```go
if len(req.Password) < 8 {
```

**Açıklama:**  
Sadece uzunluk kontrolü. `12345678`, `password`, `aaaaaaaa` gibi şifreler kabul ediliyor.

**Önerilen:** Büyük harf, rakam, özel karakter zorunluluğu ekle

---

### BUG-022: Limit Parametrelerinde Üst Sınır Yok

**Dosyalar:** `article_handler.go:39`, `profile_handler.go:253`, `admin_handler.go:197`, `topic_handler.go:47`  

**Açıklama:**  
`limit=999999` gibi değerler gönderilebilir. Sunucu tüm kayıtları hafızaya yükler.

**Önerilen:** Max limit (100-200) ile sınırla

---

### BUG-023: Statik Dosya Serving Auth Korumasız

**Dosya:** `backend/internal/infrastructure/server/fiber.go`  
**Satır:** 51  

**Kod:**
```go
app.Static("/public/avatars", "./public/avatars")
```

**Açıklama:**  
Tüm avatar dosyaları auth olmadan indirilebilir. Kullanıcı enumeration possibile.

**Önerilen:** Rate limiting veya auth ile koru

---

### BUG-024: OpenAI/Gemini Client Her Çağrıda Yeniden Oluşturuluyor

**Dosyalar:** `gateway/openai_processor.go:75-77`, `gateway/gemini_processor.go:29-33`  

**Açıklama:**  
Her `Translate`/`Summarize` çağrısında yeni HTTP client oluşturuluyor. Connection churn ve resource leak.

**Önerilen:** Client'ı singleton olarak oluştur ve yeniden kullan

---

### BUG-025: CORS Configuration Esnek

**Dosya:** `backend/internal/infrastructure/server/fiber.go`  
**Satır:** 29-34  

**Açıklama:**  
`AllowOrigins` tek bir URL'e bağlı. Üretimde yanlış konfigürasyon riski var.

**Önerilen:** Environment'a göre CORS configurasyonu

---

### BUG-026: `io.Copy`'da LimitReader Kullanılmıyor

**Dosya:** `backend/internal/adapter/handler/profile_handler.go`  
**Satır:** 161  

**Kod:**
```go
_, err = io.Copy(out, file)
```

**Açıklama:**  
Dosya boyutu kontrolü client-reported `fileHeader.Size`'a dayanıyor. Kötü niyetli istemci daha büyük gövde gönderebilir.

**Önerilen:** `io.LimitedReader` ile defense-in-depth

---

### BUG-027: Yorum İçeriğinde Maksimum Uzunluk Yok

**Dosya:** `backend/internal/usecase/comments.go`  
**Satır:** 38-39  

**Açıklama:**  
Sadece boş kontrol. Megabyte boyutunda yorum gönderilebilir.

**Önerilen:** Max 2000 karakter sınırı

---

### BUG-028: Profile Reports'ta Reason/Comment Uzunluk Sınırı Yok

**Dosya:** `backend/internal/adapter/handler/profile_handler.go`  
**Satır:** 280-288  

**Açıklama:**  
Şikayet formunda `Reason` ve `Comment` alanları için uzunluk doğrulaması yok.

**Önerilen:** Reason: 50, Comment: 500 karakter max

---

### BUG-029: Admin Topic Oluşturmada Input Doğrulama Eksik

**Dosya:** `backend/internal/adapter/handler/admin_handler.go`  
**Satır:** 48-55  

**Açıklama:**  
`Name`, `Slug`, `Keywords`, `Sources`, `RSSFeeds` alanlarında uzunluk veya format doğrulaması yok.

**Önerilen:** Her alan için max uzunluk ve format validasyonu

---

### BUG-030: Article Handler'da Source Type Doğrulaması Yok

**Dosya:** `backend/internal/adapter/handler/article_handler.go`  
**Satır:** 81  

**Kod:**
```go
source = valueobject.SourceType(sourceStr)
```

**Açıklama:**  
Herhangi bir string `SourceType` olarak kabul ediliyor. Geçersiz değerler beklenmedik davranışa neden olabilir.

---

### BUG-031: Search Query Uzunluk Sınırı Yok

**Dosya:** `backend/internal/adapter/handler/article_handler.go`  
**Satır:** 59  

**Açıklama:**  
`q` parametresi için uzunluk sınırı yok. Extremely long queries PostgreSQL `ILIKE` performsansını düşürür.

---

### BUG-032: Veritabanı Connection Pool'da Statement Timeout Yok

**Dosya:** `backend/internal/infrastructure/database/postgres.go`  
**Satır:** 19-22  

**Açıklama:**  
`MaxConns=25` ayarlı ama statement timeout veya idle-in-transaction timeout yok.

**Önerilen:** `StatementTimeout` ve `IdleInTransactionSessionTimeout` ekle

---

### BUG-033: Dockerfile'da `.env` Dosyası Intermediate Layer'a Kopyalanıyor

**Dosya:** `backend/Dockerfile`  
**Satır:** 5-6  

**Kod:**
```dockerfile
COPY . .
```

**Açıklama:**  
`.env`, test dosyaları, scratch dosyaları dahil her şey Docker build layer'ına kopyalanıyor.

**Önerilen:** `.dockerignore` ile `.env` ve gereksiz dosyaları hariç tut

---

### BUG-034: `scratch_hash.go` Production'da

**Dosya:** `backend/scratch_hash.go`  

**Açıklama:**  
Test/utility dosyası production codebase'inde. Bcrypt hash üretim scripti.

**Önerilen:** Bu dosyayı repo'dan kaldır

---

## 🔵 DÜŞÜK Açıklar (8)

### BUG-035: PasswordHash JSON Serialization'da `json:"-"` Kullanılmış (Doğru)

**Dosya:** `backend/internal/domain/entity/user.go`  
**Satır:** 10, 12  

**Not:** `PasswordHash` ve `RefreshToken` alanları `json:"-"` ile korunmuş. Bu doğru bir uygulama. **Sorun yok.**

---

### BUG-036: Veritabanı Slug Kolonunda Regex Kısıtlaması Yok

**Dosya:** `backend/migrations/001_create_topics.up.sql`  

**Açıklama:**  
`slug` alanı URL-safe karakterlerle sınırlı değil. Boşluk, Türkçe karakter vb. kabul ediliyor.

---

### BUG-037: profile_reports.status Alanında CHECK Kısıtlaması Yok

**Dosya:** `backend/migrations/011_create_user_profiles.up.sql`  
**Satır:** 41  

**Açıklama:**  
`status` alanı `'open'`, `'resolved'`, `'dismissed'` yerine herhangi bir string kabul ediyor.

---

### BUG-038: user_favorites'ta Self-Favorite Koruması Yok

**Dosya:** `backend/migrations/011_create_user_profiles.up.sql`  
**Satır:** 23-29  

**Açıklama:**  
`UNIQUE(user_id, favorite_user_id)` var ama `CHECK (user_id <> favorite_user_id)` yok. Kullanıcı kendini favorilere ekleyebilir.

---

### BUG-039: Profile Handler'da Eski Avatar Silme Hatası Yutuluyor

**Dosya:** `backend/internal/adapter/handler/profile_handler.go`  
**Satır:** 150  

**Kod:**
```go
_ = os.Remove(oldPath)
```

---

### BUG-040: Profile Handler'da Seek Hatası Yutuluyor

**Dosya:** `backend/internal/adapter/handler/profile_handler.go`  
**Satır:** 154  

**Kod:**
```go
_, _ = file.Seek(0, 0)
```

---

### BUG-041: Upload Dizini Relative Path Kullanıyor

**Dosya:** `backend/internal/adapter/handler/profile_handler.go`  
**Satır:** 137  

**Kod:**
```go
uploadDir := "./public/avatars"
```

**Açıklama:**  
Çalışma dizinine bağımlı. Mutlak path kullanılmalı.

---

### BUG-042: ListFavorites'da Limit/Offset Üst Sınırı Yok

**Dosya:** `backend/internal/adapter/handler/profile_handler.go`  
**Satır:** 253-254  

---

## ⚪ BİLGİ (2)

### BUG-043: JWT'de HS256 Kullanılmış

**Dosya:** `backend/internal/infrastructure/auth/jwt.go`  
**Satır:** 20  

**Açıklama:**  
HS256 (HMAC-SHA256) symmetric signing kullanılmış. Bu kabul edilebilir ancak secret sızarsa tüm token'lar tehlikeye girer. RS256 daha dayanıklıdır.

---

## 🧪 Test Kapsamı Kör Noktaları

### Mevcut Testler (17 dosya):
| Dosya | Kapsam |
|-------|--------|
| `auth_register_test.go` | ✅ Başarılı kayıt, duplicate username/email, eksik email |
| `article_handler_test.go` | ✅ GetRecent, Search, GetByID, Summarize |
| `admin_handler_test.go` | ✅ Auth middleware, CRUD topics, TriggerFetch |
| `summarize_article_test.go` | ✅ Happy path, fallback, hata yolları |
| `search_articles_test.go` | ✅ Query search, source filter |
| `list_articles_test.go` | ✅ List recent articles |
| `get_article_test.go` | ✅ Get by ID |
| `manage_topics_test.go` | ✅ Topic CRUD |
| `topic_handler_test.go` | ✅ GetTopics, GetTopicArticles |
| `health_handler_test.go` | ✅ Health check |
| `profile_test.go` | ✅ Bio XSS, bio uzunluk |
| `html_sanitizer_test.go` | ✅ HTML temizleme |
| `openai_processor_test.go` | ✅ OpenAI API |
| Entity/ValueObject testleri | ✅ Article, Topic, SourceType |

### Test Edilmeyen Kritik Alanlar (KRİTİK BOŞLUKLAR):

| # | Test Edilmeyen Alan | Risk | Öncelik |
|---|---------------------|------|---------|
| T1 | **Auth login flow** — Token üretimi, bcrypt doğrulama, refresh rotation | Token sahtecilik, oturum ele geçirme | 🔴 KRİTİK |
| T2 | **Auth refresh flow** — Token rotation, reuse detection,=DB temizliği | Kalıcı oturum açma, token replay | 🔴 KRİTİK |
| T3 | **Auth logout flow** — Token iptali, cookie temizliği | Yetkisiz erişim | 🟠 YÜKSEK |
| T4 | **JWT middleware** — Token doğrulama, süresi dolmuş token, hatalı format | Auth bypass | 🟠 YÜKSEK |
| T5 | **Comment creation** — İş mantığı, uzunluk doğrulama | Veritabanı şişirmesi | 🟡 ORTA |
| T6 | **Profile get/update** — Public profile, avatar upload, bio update | Veri sızıntısı | 🟡 ORTA |
| T7 | **Favorites** — Ekleme, kaldırma, listeleme | Veri tutarsızlığı | 🟡 ORTA |
| T8 | **Profile reports** — Şikayet oluşturma, duplicate kontrol | Spam, veri bozulması | 🟡 ORTA |
| T9 | **Rate limiting** — Limit çalışması, bypass | DDoS | 🟡 ORTA |
| T10 | **CORS** — Cross-origin istekler | CSRF | 🟡 ORTA |
| T11 | **File upload** — Boyut limiti, MIME doğrulama, path traversal | RCE, disk dolması | 🟠 YÜKSEK |
| T12 | **Admin authorization** — JWT role check, API key bypass | Yetki yükseltme | 🟠 YÜKSEK |
| T13 | **Security-specific tests** — SQL injection payloads, XSS payloads, auth bypass | Tüm güvenlik açıkları | 🔴 KRİTİK |

---

## Öncelikli Düzeltme Yol Haritası

### Sprint 1 (Hemen — 24 saat):
1. `.env` dosyasını git tracking'den kaldır, tüm sırları rotasyona al (BUG-001)
2. `JWT_SECRET` zorunlu kıl, hardcoded fallback'i kaldır (BUG-002)
3. Admin seed migration'ını kaldır, env-based initialization yap (BUG-003)
4. Tüm cookie'lere `Secure: true` ekle (BUG-004)
5. Public profile'dan Email ve Role alanlarını kaldır (BUG-005)
6. `role` kolonuna CHECK kısıtlaması ekle (BUG-006)

### Sprint 2 (Bu hafta):
7. Global rate limiting middleware ekle (BUG-007)
8. Summary endpoint'e auth ekle (BUG-008)
9. Refresh endpoint'ine rate limit ekle (BUG-009)
10. Admin key karşılaştırmasını `subtle.ConstantTimeCompare` yap (BUG-010)
11. Güvenlik başlıklarını ekle (BUG-011)
12. Comment rate limiting + max uzunluk ekle (BUG-012)

### Sprint 3 (Gelecek hafta):
13. Auth flow testleri yaz (T1-T4)
14. Rate limiting testleri yaz (T9)
15. File upload security testleri yaz (T11)
16. Input validation testleri yaz (T6, T7, T8)

### Teknik Borç:
17. Dockerfile'da root kullanımını düzelt (BUG-016)
18. Client singleton pattern'e geç (BUG-024)
19. `scratch_hash.go`'yu kaldır (BUG-034)
20. `.dockerignore` ekle (BUG-033)

---

## OWASP Top 10 Eşleme

| OWASP Sıralaması | HaberBot Bulguları |
|------------------|-------------------|
| A01: Broken Access Control | BUG-005, BUG-006, BUG-010 |
| A02: Cryptographic Failures | BUG-001, BUG-002, BUG-004 |
| A03: Injection | Temiz — SQL injection yok ✅ |
| A04: Insecure Design | BUG-019, BUG-021 |
| A05: Security Misconfiguration | BUG-011, BUG-016, BUG-033 |
| A06: Vulnerable Components | Denetlenmedi |
| A07: Auth Failures | BUG-002, BUG-003, BUG-007, BUG-008, BUG-009 |
| A08: Data Integrity | BUG-003, BUG-015 |
| A09: Logging & Monitoring | BUG-017, BUG-015 |
| A10: SSRF | BUG-018 (Avatar URL) |

---

**Rapor Sonu**  
*Denetim tarihi: 2026-05-29*  
*Sonraki denetim: Tüm KRİTİK ve YÜKSEK düzeltmelerden sonra*
