# Roadmap: HaberBot Webtekno Redesign

## Overview

HaberBot'un mevcut basit ızgara (grid) arayüzünü, Webtekno'nun editoryal kalitesini yansıtan zengin, çok sütunlu, reklamsız ve premium karanlık temaya sahip Webtekno görünümüne dönüştürme yol haritasıdır.

## Phases

- [x] **Phase 1: Üst Navigasyon & Kategori Menüsü** - Webtekno tarzı dinamik kategori linkleri, arama çubuğu ve sağ taraf kontrol butonları.
- [x] **Phase 2: Manşet Izgarası (Hero Grid)** - Sol sütun Sıcak Fırsatlar kartı, orta sütun büyük manşet, sağ sütun mavi/kırmızı vurgu kartı. En tepedeki boş reklam alanı elenir.
- [x] **Phase 3: Popüler Videolar Slider** - Yatayda kaydırılabilen 5 adet oynat butonlu video resim grubu.
- [x] **Phase 4: İki Sütunlu Akış & En Çok Okunanlar** - Sol sütun yatay makale feed listesi, sağ sütun numaralandırılmış trend listesi widget'ı.
- [x] **Phase 5: Kategorisel Yatay Akışlar & Alt Sitemap** - Sayfa aşağı kaydırıldıkça beliren kategori yatay listeleri, geniş site haritası linkleri ve marka iş ortaklığı footer alanı.

---

## Phase Details

### Phase 1: Üst Navigasyon & Kategori Menüsü

**Goal**: Webtekno tarzı üst navigasyon barının ve dinamik hashtag etiketlerinin HaberBot premium renk paletiyle oluşturulması.
**Depends on**: Nothing (first phase)
**Requirements**: NAV-01, NAV-02, NAV-03, NAV-04
**Success Criteria**:

  1. Kullanıcı üst barda premium "HaberBot" logosunu görür.
  2. Kullanıcı "Yapay Zeka", "Mobil" gibi kategori linklerine tıkladığında ilgili konunun haberleri React Query üzerinden tetiklenerek listelenir.
  3. Arama çubuğu, bildirimler ve profil resmi sağ tarafta Webtekno yerleşimiyle görünür.
  4. Menünün altında yatay hashtag linkleri (#Apple, #Google, vb.) listelenir.

**Plans**: 1 plan

Plans:

- [x] 01-01: Üst Navigasyon Barı ve Hashtag Etiketlerinin Tasarlanması

### Phase 2: Manşet Izgarası (Hero Grid)

**Goal**: Webtekno tarzı 3'lü manşet haber alanı ve Sıcak Fırsatlar bileşeninin reklamsız olarak ana sayfanın en tepesine yerleştirilmesi.
**Depends on**: Phase 1
**Requirements**: GRID-01, GRID-02, GRID-03, GRID-04
**Success Criteria**:

  1. Sayfanın en üstünde reklam boşluğu olmadan doğrudan haberler başlar.
  2. Sol tarafta mor renkli "Sıcak Fırsatlar & İndirimler" kartı dikeyde listelenir.
  3. Ortada büyük manşet resmi ve yeşil-neon bindirme alt başlıklı haber kartı yer alır.
  4. Sağ tarafta düz elektrik-mavisi arka planlı, sade ve beyaz metinli çarpıcı vurgu haber kartı yer alır.

**Plans**: 1 plan

Plans:

- [x] 02-01: 3 Sütunlu Kahraman Izgarasının (Sıcak Fırsatlar, Manşet, Mavi Vurgu) Kodlanması

### Phase 3: Popüler Videolar Slider

**Goal**: Ana manşet alanının altında, yatayda kaydırılabilen "Popüler Videolar" slider bileşeninin oluşturulması.
**Depends on**: Phase 2
**Requirements**: SLID-01, SLID-02
**Success Criteria**:

  1. Kırmızı oynat ikonlu "Popüler Videolar" başlığı eklenmiştir.
  2. Başlığın altında, yatayda scroll edilebilen, üzerinde oynat butonu bulunan 5 adet video küçük resmi yer alır.

**Plans**: 1 plan

Plans:

- [x] 03-01: Popüler Videolar Slider Bileşeninin Entegre Edilmesi

### Phase 4: İki Sütunlu Akış & En Çok Okunanlar

**Goal**: Sol tarafta yatay makale listesi ve sağ tarafta "En Çok Okunanlar" trend listesini barındıran iki sütunlu akışın oluşturulması.
**Depends on**: Phase 3
**Requirements**: FEED-01, FEED-02
**Success Criteria**:

  1. Sol tarafta yatay yerleşimli büyük görselli makale kartları dikey olarak listelenir.
  2. Sağ tarafta "EN ÇOK OKUNANLAR" (Most Read) numaralandırılmış trend listesi widget'ı doğru şekilde render edilir.

**Plans**: 1 plan

Plans:

- [x] 04-01: İki Sütunlu Akış (Sol Makale Listesi, Sağ Trend Sidebar) Oluşturulması

### Phase 5: Kategorisel Yatay Akışlar & Alt Sitemap

**Goal**: "Otomobil" gibi kategorilere özel yatay akış grupları, site haritası sitemap linkleri ve footer alanının yerleştirilmesi.
**Depends on**: Phase 4
**Requirements**: SITE-01, SITE-02, SITE-03
**Success Criteria**:

  1. Sayfa aşağı kaydırıldıkça "Otomobil" başlığı ve yatay kaydırılabilir kart grupları listelenir.
  2. Footer'da "Kurumsal", "Kategoriler" sitemap linkleri ve partner marka logoları yer alır.

**Plans**: 1 plan

Plans:

- [x] 05-01: Kategorisel Yatay Akışların ve Geniş Sitemap Footer Alanının Eklenmesi

### Phase 6: Pending queue limit for onay bekleyen haberler

**Goal:** Onay bekleyen haberlerde, kategori başına pending sayısını config ile belirlenen limitte tutmak (varsayılan 50).
**Depends on:** Nothing (independent backend change)
**Requirements**: PEND-01, PEND-02, PEND-03, PEND-04
**Status:** Closed (user requested)
**Success Criteria**:

  1. Pending haberlerde kategori başına limit uygulanır ve en eski `fetched_at` kayıtları silinir.
  2. Limit yalnızca pending (is_approved=false) haberler için geçerlidir; onaylılar etkilenmez.
  3. Limit değeri `PENDING_LIMIT_PER_TOPIC` ile değiştirilebilir; boşsa 50 kullanılır.

**Plans:** 0 plans

Plans:

- [x] Closed — no implementation planned

### Phase 7: Auth & Role Management

**Goal:** Üyelik, rol yönetimi ve yorum altyapısını güvenli auth akışıyla devreye almak.
**Requirements**: AUTH-01, AUTH-02, AUTH-03, AUTH-04, AUTH-05
**Depends on:** Nothing (independent backend change)
**Plans:** 1 plan

Plans:

- [x] 07-01: Üyelik, Rol Yönetimi ve Yorum Altyapısının Uçtan Uca Entegrasyonu (Tabbed AuthModal, E-Posta doğrulaması ve birim testler tamamlandı!)

### Phase 8: User Profile

**Goal:** Kullanıcı profil sayfası, profil düzenleme (bio + avatar), favorilere ekleme ve profil raporlama sistemini devreye almak.
**Depends on:** Phase 7 (Auth & Role Management)
**Requirements**: PROF-01, PROF-02, PROF-03, PROF-04, PROF-05, PROF-06, PROF-07, PROF-08, PROF-09
**Status:** Completed
**Plans:** 1 plan

Plans:

- [x] 08-01: Kullanıcı Profil Sayfası, Avatar Yönetimi, Favoriler ve Raporlama Sistemi (18 görevli kapsamlı uygulama planı)

### Phase 9: Smart TL;DR & AI Summary Redesign

**Goal:** Haber detayından gömülü AI özeti kaldırıp sticky CTA + premium modal akışıyla lazy-loaded, markdown-safe özet deneyimi sunmak.
**Depends on:** Phase 8 (User Profile)
**Requirements**: SUM-01, SUM-02, SUM-03, SUM-04, SUM-05
**Status:** Planning
**Plans:** 1 plan

Plans:

- [ ] 09-01: Sticky AI Summary CTA, Modal, Markdown Render ve Lazy Summary Endpoint Entegrasyonu

---

## Progress

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Üst Navigasyon | 1/1 | Completed | 2026-05-29 |
| 2. Manşet Izgarası | 1/1 | Completed | 2026-05-29 |
| 3. Popüler Videolar | 1/1 | Completed | 2026-05-29 |
| 4. İki Sütunlu Akış | 1/1 | Completed | 2026-05-29 |
| 5. Kategorisel Akışlar | 1/1 | Completed | 2026-05-29 |
| 6. Pending queue limit | 0/0 | Closed | 2026-05-29 |
| 7. Auth & Role Management | 1/1 | Completed | 2026-05-29 |
| 8. User Profile | 1/1 | Completed | 2026-05-29 |
| 9. Smart TL;DR & AI Summary | 0/1 | Planning | — |
