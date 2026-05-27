# HaberBot Webtekno Redesign

## What This Is

HaberBot, modern teknoloji ve genel haber akışlarını toplayan, yapay zeka ile özetleyip Türkçe editoryal kalitede sunan bir haber platformudur. Bu projenin amacı, HaberBot web arayüzünün (homepage ve ana menüleri) Türkiye'nin lider teknoloji portalı **Webtekno**'nun görsel yerleşimi ve ızgara (grid) yapısıyla birebir eşleşecek şekilde tamamen yeniden tasarlanmasıdır.

## Core Value

HaberBot'un premium karanlık mod renk paletini ("Midnight Neon") koruyarak, Webtekno'nun zengin, dinamik ve çok sütunlu editoryal yerleşimini hatasız ve yüksek performansla son kullanıcıya sunmak.

## Requirements

### Validated

- ✓ Go Clean Architecture (Ports & Adapters) backend platformu — Phase 0
- ✓ PostgreSQL veritabanı şeması ve pgxpool entegrasyonu — Phase 0
- ✓ React 19 ve TanStack React Query v5 veri çekme ve durum önbellek altyapısı — Phase 0
- ✓ Gemini API / OpenAI API ile Türkçe çeviri ve özetleme mekanizması — Phase 0

### Active

- [ ] **Webtekno Tarzı Üst Navigasyon & Kategori Menüsü:** Webtekno'daki gibi sol tarafta logo, orta kısımda kategoriler ("Yapay Zeka", "Mobil", vb.), sağ tarafta ise arama, profil ve karanlık/aydınlık mod anahtarı barındıran menü barı.
- [ ] **Reklamsız Üst Alan Yerleşimi:** Orijinal Webtekno sitesinde en üstte yer alan boş reklam baneri alanı tamamen iptal edilecek, içerikler ve "Popüler İçerikler" etiketleri doğrudan navigasyonun hemen altından başlayacaktır.
- [ ] **Webtekno Tarzı Kahraman (Hero) Izgarası (Grid):** Sol tarafta "Sıcak Fırsatlar" kartı (özelleştirilmiş kampanya kartı), ortada büyük görsel manşet haber kartı, sağda ise kırmızı arka planlı vurgulu başlık kartı.
- [ ] **Popüler Videolar Kaydırıcısı (Slider):** Ana ızgaranın hemen altında konumlandırılmış, yatayda kaydırılabilen (horizontal slider) popüler videolar bileşeni.
- [ ] **İki Sütunlu Alt İçerik Akışı:** Sol sütunda yatay yerleşimli büyük görselli makale listesi, sağ sütunda ise "En Çok Okunanlar" (Most Read) numaralandırılmış liste widget'ı.
- [ ] **Webtekno Tarzı Kategorisel Yatay Akışlar:** Sayfa aşağı kaydırıldıkça "Otomobil", "Uygulama/Yazılım" gibi kategorilere özel, kaydırılabilir yatay kart grupları ve listeleri.
- [ ] **Kapsamlı Footer (Site Haritası):** Webtekno'daki gibi kategorize edilmiş detaylı sitemap linkleri, sosyal medya ikonları ve marka iş ortaklığı logoları.

### Out of Scope

- **Reklam Yerleşimleri:** Webtekno'daki reklam alanlarının projeye eklenmesi — Kullanıcı talebi doğrultusunda reklam baneri alanları tamamen elenerek haberlerin en tepeden başlaması sağlanacaktır.
- **Renk Paleti Değişikliği:** Webtekno'nun beyaz/kırmızı renk şemasına geçiş — Projenin mevcut premium karanlık mod estetiği ("Midnight Neon" ve zifiri siyah zemin) korunacaktır.

## Context

- **Mevcut Durum:** HaberBot şu anda tek sütunlu / basit 3 sütunlu standart bir kart ızgarasına sahiptir. Bu tasarım, çok sayıda haber geldiğinde editoryal zenginliği tam olarak yansıtamamaktadır.
- **Görsel Standartlar:** `DESIGN.md` dosyasındaki karanlık mod prensipleri (`#000000` zemin, `#3b82f6` mavi vurgu vb.) ve ince vuruşlu kenarlıklar (`1px solid #222222`) yeni arayüze birebir uyarlanacaktır.

## Constraints

- **Tasarım:** Arayüz bileşenlerinin yerleşimi Webtekno'nun birebir kopyası olmalıdır.
- **Renkler:** Orijinal sitenin renk paleti yerine HaberBot'un kendi koyu/neon paleti kullanılacaktır.
- **Performans:** Resim yüklemeleri ve kaydırmalar pürüzsüz olmalı, TanStack Query önbelleği bozulmamalıdır.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Webtekno Izgarası (Grid) ve Slider yapısını entegre etmek | Editoryal magazin kalitesini artırmak ve zengin içerik sunumu sağlamak | — Pending |
| Koyu temayı (`#000000`) korumak | Kullanıcı marka kimliğini ve premium dijital atmosferi sürdürmek | — Pending |
| Boş reklam alanını (Ad banner) kaldırmak | Sayfa yükleme hızını artırmak ve gereksiz boşlukları önlemek | — Pending |

---
*Last updated: 2026-05-27 after initialization*
