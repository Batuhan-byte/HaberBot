# Roadmap: HaberBot Webtekno Redesign

## Overview

HaberBot'un mevcut basit ızgara (grid) arayüzünü, Webtekno'nun editoryal kalitesini yansıtan zengin, çok sütunlu, reklamsız ve premium karanlık temaya sahip Webtekno görünümüne dönüştürme yol haritasıdır.

## Phases

- [ ] **Phase 1: Üst Navigasyon & Kategori Menüsü** - Webtekno tarzı dinamik kategori linkleri, arama çubuğu ve sağ taraf kontrol butonları.
- [ ] **Phase 2: Manşet Izgarası (Hero Grid)** - Sol sütun Sıcak Fırsatlar kartı, orta sütun büyük manşet, sağ sütun mavi/kırmızı vurgu kartı. En tepedeki boş reklam alanı elenir.
- [ ] **Phase 3: Popüler Videolar Slider** - Yatayda kaydırılabilen 5 adet oynat butonlu video resim grubu.
- [ ] **Phase 4: İki Sütunlu Akış & En Çok Okunanlar** - Sol sütun yatay makale feed listesi, sağ sütun numaralandırılmış trend listesi widget'ı.
- [ ] **Phase 5: Kategorisel Yatay Akışlar & Alt Sitemap** - Sayfa aşağı kaydırıldıkça beliren kategori yatay listeleri, geniş site haritası linkleri ve marka iş ortaklığı footer alanı.

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
- [ ] 01-01: Üst Navigasyon Barı ve Hashtag Etiketlerinin Tasarlanması

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
- [ ] 02-01: 3 Sütunlu Kahraman Izgarasının (Sıcak Fırsatlar, Manşet, Mavi Vurgu) Kodlanması

### Phase 3: Popüler Videolar Slider
**Goal**: Ana manşet alanının altında, yatayda kaydırılabilen "Popüler Videolar" slider bileşeninin oluşturulması.
**Depends on**: Phase 2
**Requirements**: SLID-01, SLID-02
**Success Criteria**:
  1. Kırmızı oynat ikonlu "Popüler Videolar" başlığı eklenmiştir.
  2. Başlığın altında, yatayda scroll edilebilen, üzerinde oynat butonu bulunan 5 adet video küçük resmi yer alır.
**Plans**: 1 plan

Plans:
- [ ] 03-01: Popüler Videolar Slider Bileşeninin Entegre Edilmesi

### Phase 4: İki Sütunlu Akış & En Çok Okunanlar
**Goal**: Sol tarafta yatay makale listesi ve sağ tarafta "En Çok Okunanlar" trend listesini barındıran iki sütunlu akışın oluşturulması.
**Depends on**: Phase 3
**Requirements**: FEED-01, FEED-02
**Success Criteria**:
  1. Sol tarafta yatay yerleşimli büyük görselli makale kartları dikey olarak listelenir.
  2. Sağ tarafta "EN ÇOK OKUNANLAR" (Most Read) numaralandırılmış trend listesi widget'ı doğru şekilde render edilir.
**Plans**: 1 plan

Plans:
- [ ] 04-01: İki Sütunlu Akış (Sol Makale Listesi, Sağ Trend Sidebar) Oluşturulması

### Phase 5: Kategorisel Yatay Akışlar & Alt Sitemap
**Goal**: "Otomobil" gibi kategorilere özel yatay akış grupları, site haritası sitemap linkleri ve footer alanının yerleştirilmesi.
**Depends on**: Phase 4
**Requirements**: SITE-01, SITE-02, SITE-03
**Success Criteria**:
  1. Sayfa aşağı kaydırıldıkça "Otomobil" başlığı ve yatay kaydırılabilir kart grupları listelenir.
  2. Footer'da "Kurumsal", "Kategoriler" sitemap linkleri ve partner marka logoları yer alır.
**Plans**: 1 plan

Plans:
- [ ] 05-01: Kategorisel Yatay Akışların ve Geniş Sitemap Footer Alanının Eklenmesi

---

## Progress

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Üst Navigasyon | 0/1 | Not started | - |
| 2. Manşet Izgarası | 0/1 | Not started | - |
| 3. Popüler Videolar | 0/1 | Not started | - |
| 4. İki Sütunlu Akış | 0/1 | Not started | - |
| 5. Kategorisel Akışlar | 0/1 | Not started | - |
