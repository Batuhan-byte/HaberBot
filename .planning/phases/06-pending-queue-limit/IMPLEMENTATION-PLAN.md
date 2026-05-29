# Implementasyon Planı: Pending Queue Limit

**Tarih:** 2026-05-29  
**Versiyon:** 1.0  
**Durum:** Planlandı

---

## Özet

Onay bekleyen haberlerde kategori başına maksimum pending sayısını (configurable, varsayılan 50) korumak için insert sonrası trim ve periyodik temizlik akışları eklenecek. Silme kriteri `fetched_at` (en eski önce) olacak ve sadece pending kayıtlar etkilenecek.

---

## Mevcut Durum Analizi

- Pending haberlerde limit uygulanmıyor.
- Periyodik temizlik job’ı pending limit için yok.
- `topic_id` NULL olan haberler için özel limit davranışı tanımlı değil.

---

## Faz 01: Veri Katmanı (Repository + Config)

### 1.1 Config’e limit alanı ekleme
**Hedef:** `PENDING_LIMIT_PER_TOPIC` env desteği ve varsayılan 50.

**Dosyalar:**
- `backend/internal/infrastructure/config/config.go`

---

### 1.2 Repository trim metodu
**Hedef:** `TrimPendingByTopic(topicID, limit)` metodu ile en eski pending kayıtları silmek.

**Dosyalar:**
- `backend/internal/domain/port/article_repository.go`
- `backend/internal/adapter/repository/postgres_article.go`

---

## Faz 02: Usecase Akışı

### 2.1 Insert sonrası trim
**Hedef:** Yeni pending haber eklendiğinde ilgili topic için trim çalıştırmak.

**Dosyalar:**
- `backend/internal/usecase/fetch_articles.go` (veya pending insert yapılan akış)

---

### 2.2 Periyodik temizlik
**Hedef:** Scheduler üzerinden tüm topic’ler (NULL dahil) için trim çalıştırmak.

**Dosyalar:**
- `backend/internal/infrastructure/scheduler/scheduler.go`
- `backend/internal/usecase/daily_pipeline.go` (varsa uygun akış)

---

## Faz 03: Test ve Doğrulama

### 3.1 Repository testleri
**Hedef:** Limit aşıldığında en eski pending kayıtların silindiğini doğrulamak.

**Dosyalar:**
- `backend/internal/adapter/repository/postgres_article_test.go`

---

### 3.2 Usecase testleri
**Hedef:** Insert sonrası trim akışı ve NULL topic senaryosu.

**Dosyalar:**
- `backend/internal/usecase/..._test.go`

---

## Görev Listesi

| # | Görev | Dosyalar | Durum |
|---|-------|----------|-------|
| 1 | Config’e limit alanı | config.go | Beklemede |
| 2 | Repo trim metodu | article_repository.go, postgres_article.go | Beklemede |
| 3 | Insert sonrası trim | fetch_articles.go | Beklemede |
| 4 | Scheduler trim | scheduler.go, daily_pipeline.go | Beklemede |
| 5 | Testler | *_test.go | Beklemede |

---

## Success Criteria

- [ ] Pending haberler kategori başına limite çekiliyor.
- [ ] En eski `fetched_at` kayıtları siliniyor.
- [ ] Sadece pending haberler etkileniyor; approved haberler korunuyor.
- [ ] NULL topic bucket’ı ayrı limitleniyor.
- [ ] Config ile limit değiştirilebiliyor.

---

## Bağımlılıklar

1 → 2 → 3 → 4 → 5

