# Phase 6: Pending Queue Limit per Topic — Specification

**Created:** 2026-05-29  
**Ambiguity score:** 0.14 (gate: ≤ 0.20)  
**Requirements:** 6 locked

## Goal

Onay bekleyen haberlerde, kategori başına pending sayısını config ile belirlenen limitte (varsayılan 50) tutmak ve limit aşıldığında en eski `fetched_at` kayıtlarını silmek.

## Background

Pending haber kuyruğu sınırsız büyüyor ve bazı kategorilerde yüzlerce pending kayıt birikiyor. Mevcut akışta pending insert sonrası kırpma yapılmıyor, ayrıca periyodik temizlik yok. Bu, gereksiz depolama maliyeti ve yönetim zorluğu yaratıyor.

## Requirements

1. **Configurable limit**: Pending limit değeri config üzerinden okunacak.
   - Current: Böyle bir limit yok; config’te ilgili bir alan yok.
   - Target: `PENDING_LIMIT_PER_TOPIC` env varı ile limit belirlenebilecek, boşsa varsayılan 50 kullanılacak.
   - Acceptance: Env varı set edilince limit değeri değişiyor; set değilken 50 kullanılıyor.

2. **Anlık kırpma (insert sonrası)**: Yeni pending haber eklendiğinde ilgili topic için limit aşımı varsa en eskiler silinecek.
   - Current: Insert sonrası pending limit kontrolü yok.
   - Target: Pending insert akışının sonunda `TrimPendingByTopic(topic_id, limit)` çağrılır ve limit üstü pending kayıtlar silinir.
   - Acceptance: Tek topic altında 51. pending eklendiğinde toplam pending sayısı 50’ye düşer.

3. **Periyodik temizlik**: Scheduler düzenli şekilde tüm topic’ler için limit kontrolü yapacak.
   - Current: Pending temizlik için periyodik job yok.
   - Target: Cron/scheduler, tüm topic_id (NULL dahil) için trim uygular.
   - Acceptance: Scheduler tetiklendiğinde her topic’te pending sayısı limite çekilir.

4. **Sadece pending**: Limit uygulaması sadece onay bekleyen (is_approved = false) haberleri kapsar.
   - Current: Pending/approved ayrımı yapılmadan limit yok.
   - Target: Approved haberler asla silinmez; yalnızca pending kırpılır.
   - Acceptance: Approved haberler limit koşulunda etkilenmez.

5. **NULL topic davranışı**: `topic_id` NULL olan pending haberler ayrı bir kategori gibi değerlendirilir.
   - Current: NULL topic için davranış tanımlı değil.
   - Target: NULL topic bağımsız bir bucket olarak limitlenir.
   - Acceptance: NULL topic pending sayısı limitten fazla ise en eskiler silinir.

6. **Hard delete**: Silme işlemi kalıcı (DELETE) olacak.
   - Current: Pending haberler için otomatik silme yok.
   - Target: Limit üstü pending haberler DB’den hard delete edilir.
   - Acceptance: Trim sonrası silinen kayıtlar DB’de bulunamaz.

## Boundaries

**In scope:**
- Config’e `PENDING_LIMIT_PER_TOPIC` eklenmesi ve varsayılan 50.
- Repository seviyesinde `TrimPendingByTopic` (topic_id + pending) kırpma işlemi.
- Pending insert akışında trim çağrısı.
- Scheduler/cron ile periyodik trim.

**Out of scope:**
- Onaylı haberler için limit/retention politikası — sadece pending içindir.
- Frontend UI değişiklikleri — mevcut ekranlar aynı kalır.
- Soft delete/arşivleme — yalnızca hard delete.

## Constraints

- Eski/yeniyi belirleme kriteri **fetched_at** olacak.
- Pending limit mantığı usecase → repository katmanlarında uygulanır; handler/route değişikliği gerekmez.

## Acceptance Criteria

- [ ] `PENDING_LIMIT_PER_TOPIC` env varı set edilince limit buna göre uygulanır; set değilse 50 kullanılır.
- [ ] Aynı topic için 51. pending eklendiğinde pending sayısı 50’ye düşer.
- [ ] Approved (is_approved = true) haberler hiçbir koşulda silinmez.
- [ ] NULL topic bucket’ı da limitlenir.
- [ ] Trim işlemi en eski `fetched_at` kayıtlarını siler.
- [ ] Scheduler tetiklendiğinde tüm topic’lerde pending sayısı limite çekilir.

## Ambiguity Report

| Dimension           | Score | Min  | Status | Notes                             |
|--------------------|-------|------|--------|-----------------------------------|
| Goal Clarity        | 0.92  | 0.75 | ✓      | Net hedef ve ölçüm                |
| Boundary Clarity    | 0.90  | 0.70 | ✓      | In/Out scope belirgin             |
| Constraint Clarity  | 0.78  | 0.65 | ✓      | fetched_at kriteri net            |
| Acceptance Criteria | 0.82  | 0.70 | ✓      | Tamamı test edilebilir            |
| **Ambiguity**       | 0.14  | ≤0.20| ✓      |                                   |

## Interview Log

| Round | Perspective     | Question summary                                  | Decision locked                                               |
|-------|-----------------|--------------------------------------------------|---------------------------------------------------------------|
| 1     | Researcher      | Limit pending mi, toplam mı?                     | Sadece pending (onay bekleyen)                                |
| 2     | Simplifier      | Eskiyi ne belirlesin?                             | fetched_at                                                     |
| 3     | Boundary Keeper | Silme şekli?                                      | Hard delete                                                    |
| 4     | Failure Analyst | Ne zaman çalışsın?                                | Insert sonrası + periyodik scheduler                           |
| 5     | Boundary Keeper | NULL topic nasıl davranır?                        | Ayrı kategori gibi limitlenir                                  |
| 6     | Simplifier      | Limit sabit mi?                                   | Config üzerinden değişebilir (varsayılan 50)                   |

---

*Phase: 06-pending-queue-limit*  
*Spec created: 2026-05-29*  
*Next step: /gsd-plan-phase 6 — implementation plan*
