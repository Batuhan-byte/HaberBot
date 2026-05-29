# Phase 7: Auth & Role Management — Specification

**Created:** 2026-05-29  
**Ambiguity score:** 0.12 (gate: ≤ 0.20)  
**Requirements:** 9 locked

## Goal

HaberBot’ta üyelik, rol yönetimi ve yorum altyapısını güvenli auth akışıyla devreye almak; Admin ve User rollerini ayırmak, yorum yazmayı User rolüne bağlamak ve admin panelini role-based auth ile korumak.

## Background

Mevcut sistemde admin paneli için API key tabanlı koruma mevcut; kullanıcı kayıt/giriş ve rol tabanlı yetkilendirme yok. Yorum özelliği bulunmuyor. Yeni tasarım, modern ve güvenli auth akışıyla üyelik, rol ve yorum altyapısını eklemeyi hedefliyor.

## Requirements

1. **Rol modeli**: Sistem yalnızca `Admin` ve `User` rollerini destekler.
   - Current: Rol yapısı yok.
   - Target: `users` tablosunda rol alanı eklenir ve yetkilendirme bu role göre yapılır.
   - Acceptance: Admin-only endpoint User rolü ile erişilemez.

2. **Auth yöntemi**: Kısa ömürlü JWT access token + HttpOnly refresh cookie.
   - Current: JWT/cookie akışı yok.
   - Target: Access token response ile döner; refresh token HttpOnly+Secure+SameSite cookie olarak set edilir.
   - Acceptance: Access token 15–30 dk ömürlü, refresh token ile yenilenebilir.

3. **Kayıt/Giriş/Çıkış API’leri**: Auth endpoint’leri eklenecek.
   - Current: Auth API yok.
   - Target: `POST /api/auth/register`, `POST /api/auth/login`, `POST /api/auth/refresh`, `POST /api/auth/logout`.
   - Acceptance: Başarılı login sonrası access token döner, refresh cookie set edilir; logout cookie temizler.

4. **Şifre güvenliği**: Şifreler bcrypt ile hash’lenir.
   - Current: Kullanıcı şifre alanı yok.
   - Target: Kullanıcı şifresi hiçbir zaman düz metin saklanmaz; bcrypt hash DB’de tutulur.
   - Acceptance: DB’de düz metin şifre bulunmaz.

5. **Seed admin hesabı**: İlk kurulumda `admin1/admin1` Admin olarak yaratılır.
   - Current: Seed admin yok.
   - Target: Migration/seed akışında admin1 kullanıcısı bcrypt hash ile oluşturulur.
   - Acceptance: İlk çalıştırmada Admin rolünde bir kullanıcı oluşur.

6. **Admin paneli koruması**: Role-based auth + mevcut API key birlikte.
   - Current: Sadece API key koruması var.
   - Target: Admin UI endpoint’leri role-based auth kontrolü yapar; API key kontrolü internal işlemler için ayrı kalır.
   - Acceptance: Admin rolü olmayan kullanıcılar panel endpoint’lerine erişemez.

7. **Yorum altyapısı**: Yeni `comments` tablosu ve API.
   - Current: Yorum tablosu yok.
   - Target: `comments` tablosu + `POST /api/comments`, `GET /api/comments?article_id=...`.
   - Acceptance: User rolü yorum oluşturabilir; yorumlar public listelenir.

8. **Validation**: Username uniqueness ve minimum şifre uzunluğu (8).
   - Current: Validasyon yok.
   - Target: Username benzersiz, şifre >= 8.
   - Acceptance: Aynı username ile kayıt engellenir; kısa şifre 400 döner.

9. **Güvenlik önlemleri**: Rate limit (login/register), refresh token rotation.
   - Current: Rate limit ve rotation yok.
   - Target: Login/register için rate limit uygulanır; refresh token kullanıldığında yenisi üretilir.
   - Acceptance: Çoklu deneme sınırlandı; refresh sonrası eski token geçersiz olur.

## Boundaries

**In scope:**
- Auth API ve JWT/cookie akışı
- Users ve Comments tabloları
- Role-based auth middleware
- Seed admin oluşturma

**Out of scope:**
- OAuth/social login — ayrı faz
- Email doğrulama — ayrı faz
- Admin UI tasarım değişiklikleri — mevcut arayüz korunur

## Constraints

- Şifre hashleme bcrypt ile yapılacak.
- Refresh cookie HttpOnly + Secure + SameSite olacak.
- Yorumlar admin onayı olmadan yayınlanır.

## Acceptance Criteria

- [ ] Admin ve User rollerine göre erişim kontrolü çalışır.
- [ ] Login → access token döner; refresh cookie set edilir.
- [ ] Refresh endpoint yeni access token üretir, refresh token rotate eder.
- [ ] Şifreler bcrypt hash’li saklanır; düz metin saklanmaz.
- [ ] `admin1/admin1` seed admin oluşur.
- [ ] `POST /api/comments` yalnızca User ile çalışır; list endpoint public çalışır.
- [ ] Username benzersizliği ve min 8 şifre zorunluluğu uygulanır.
- [ ] Login/register rate limit çalışır.

## Ambiguity Report

| Dimension           | Score | Min  | Status | Notes                             |
|--------------------|-------|------|--------|-----------------------------------|
| Goal Clarity        | 0.92  | 0.75 | ✓      | Net hedef                         |
| Boundary Clarity    | 0.88  | 0.70 | ✓      | In/Out scope açık                 |
| Constraint Clarity  | 0.80  | 0.65 | ✓      | Security constraints net          |
| Acceptance Criteria | 0.84  | 0.70 | ✓      | Test edilebilir kriterler         |
| **Ambiguity**       | 0.12  | ≤0.20| ✓      |                                   |

## Interview Log

| Round | Perspective     | Question summary                                  | Decision locked                                               |
|-------|-----------------|--------------------------------------------------|---------------------------------------------------------------|
| 1     | Boundary Keeper | Admin API key korunumu nasıl olacak?             | Role-based auth + API key ayrı kalsın                         |
| 2     | Simplifier      | Kullanıcı kimliği alanları?                      | Username + password                                           |
| 3     | Boundary Keeper | Yorumlar onaylı mı?                               | Hemen görünür                                                 |
| 4     | Security        | Token modeli?                                    | Short-lived access + refresh cookie                           |
| 5     | Security        | Min şifre uzunluğu?                              | 8                                                             |

---

*Phase: 07-auth-role-management*  
*Spec created: 2026-05-29*  
*Next step: /gsd-plan-phase 7 — implementation plan*
