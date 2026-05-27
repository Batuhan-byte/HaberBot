# Requirements: HaberBot Webtekno Redesign

**Defined:** 2026-05-27
**Core Value:** HaberBot'un premium karanlık mod renk paletini ("Midnight Neon") koruyarak, Webtekno'nun zengin, dinamik ve çok sütunlu editoryal yerleşimini hatasız ve yüksek performansla son kullanıcıya sunmak.

## v1 Requirements

### 1. Üst Navigasyon & Kategori Menüsü (Nav & Menu)

- [ ] **NAV-01:** Sol tarafta premium "HaberBot" logosu yer almalıdır.
- [ ] **NAV-02:** Orta kısımda yatayda dizilmiş dinamik kategori menü butonları ("Yapay Zeka", "Mobil", "Ürün", "Oyun", "Otomobil", "Uygulama/Yazılım", "Donanım", "Sinema-Dizi") bulunmalı ve tıklandığında ilgili konunun (Topic) haberleri yüklenmelidir.
- [ ] **NAV-03:** Sağ tarafta arama çubuğu (search), bildirim ikonu (notifications), yönetici ayarları (settings) ve Unsplash profil resmi ile giriş butonu yer almalıdır.
- [ ] **NAV-04:** Kategori menüsünün hemen altında "Popüler İçerikler" etiketi ve yanında yatayda dizilmiş popüler hashtag linkleri (#Apple, #Google, vb.) bulunmalıdır.

### 2. Kahraman (Hero) Izgarası & Reklamsız Üst Alan (Hero Grid)

- [ ] **GRID-01:** Webtekno'daki üst boş reklam baneri alanı tamamen iptal edilmelidir; "Popüler İçerikler" etiketlerinin hemen altında doğrudan ana haber ızgarası başlamalıdır.
- [ ] **GRID-02:** Sol Sütun ("Sıcak Fırsatlar & İndirimler"): Mor arka planlı, dikey bir kart bulunmalıdır. İçerisinde editoryal fırsat listeleri, indirim kuponları ve "webtekno30" tarzı indirim kodları listelenmelidir.
- [ ] **GRID-03:** Orta Sütun (Büyük Manşet Kartı): En güncel veya en yüksek puanlı (score) haberin büyük manşet resmi, görsel üzerine bindirilmiş başlığı ve yeşil-neon renkli "NASIL DEĞİŞTİRİLİR?" veya benzeri dikkat çekici alt başlığı yer almalıdır.
- [ ] **GRID-04:** Sağ Sütun (Vurgulu Başlık Kartı): Kırmızı (`#ef4444`) veya projenin mavi tonlu arka planına sahip, düz renkli ve üzerinde sadece beyaz renkli çarpıcı başlık ve kısa açıklama barındıran vurgu haber kartı yer almalıdır.

### 3. Popüler Videolar Slider Bileşeni (Video Slider)

- [ ] **SLID-01:** Ana manşet ızgarasının altında, kırmızı renkli bir oynatma ikonuyla "Popüler Videolar" başlığı bulunmalıdır.
- [ ] **SLID-02:** Başlığın altında yatayda kaydırılabilen (horizontal scrolling), üzerinde oynat düğmesi bulunan ve tıklandığında ilgili video/detay sayfasını açan 5 adet video küçük resmi yer almalıdır.

### 4. İki Sütunlu Alt İçerik Akışı (Feed & Sidebar)

- [ ] **FEED-01:** Sol taraftaki geniş sütunda, yatay yerleşimli büyük görselli makale kartları dikey olarak listelenmelidir. Her kartta sağda veya solda büyük görsel, yanında başlık, kısa özet, yazar ismi ve "3 saat önce" gibi zaman damgası yer almalıdır.
- [ ] **FEED-02:** Sağ taraftaki dar sütunda "EN ÇOK OKUNANLAR" (Most Read) widget'ı bulunmalıdır. Widget'ın en üstünde popüler bir büyük haber kartı, altında ise 1'den başlayıp numaralandırılmış trend haber başlıkları listelenmelidir.

### 5. Kategorisel Yatay Akışlar & Alt Sitemap (Sitemap)

- [ ] **SITE-01:** Sayfa aşağı kaydırıldıkça "Otomobil" gibi kategorilere ait özel yatay kaydırılabilir kart grupları ikonlarıyla birlikte listelenmelidir.
- [ ] **SITE-02:** Sayfanın en altında Webtekno tarzı geniş site haritası (sitemap) linkleri ("Kurumsal", "Kategoriler", "En Çok Paylaşılanlar", vb.) yer almalıdır.
- [ ] **SITE-03:** Footer'ın en altında gri sosyal medya ikonları, "Bir mediaone markasıdır." alt yazısı ve mackolik, onedio tarzı partner marka logoları bulunmalıdır.

## v2 Requirements

### Arama & Filtreler (Deferred)

- **SRCH-01:** Arama giriş çubuğunda anlık arama (instant autocomplete) önerileri listelenmesi.
- **DARK-01:** Webtekno tarzı ışık/karanlık mod geçiş anahtarının (light/dark switch) işlevsel hale getirilmesi (şu an sadece görsel tema korunacaktır).

## Out of Scope

| Feature | Reason |
|---------|--------|
| Reklam Alanları (Ad Banners) | Kullanıcı talebi üzerine reklamsız, temiz yerleşim hedeflenmiştir. |
| Orijinal Webtekno Renk Paleti | Projenin kendi premium "Midnight Neon" koyu teması korunacaktır. |

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| NAV-01 | Phase 1 | Pending |
| NAV-02 | Phase 1 | Pending |
| NAV-03 | Phase 1 | Pending |
| NAV-04 | Phase 1 | Pending |
| GRID-01 | Phase 2 | Pending |
| GRID-02 | Phase 2 | Pending |
| GRID-03 | Phase 2 | Pending |
| GRID-04 | Phase 2 | Pending |
| SLID-01 | Phase 3 | Pending |
| SLID-02 | Phase 3 | Pending |
| FEED-01 | Phase 4 | Pending |
| FEED-02 | Phase 4 | Pending |
| SITE-01 | Phase 5 | Pending |
| SITE-02 | Phase 5 | Pending |
| SITE-03 | Phase 5 | Pending |

**Coverage:**
- v1 requirements: 15 total
- Mapped to phases: 15
- Unmapped: 0 ✓

---
*Requirements defined: 2026-05-27*
*Last updated: 2026-05-27 after initial definition*
