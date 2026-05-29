# Tasarım Özellikleri: Dinamik RSS ve Akıllı Filtreleme Haber Motoru Revizyonu

Bu döküman, HaberBot projesinin haber çekme (Fetch) motorunu canlandırmak, zenginleştirmek, kategorilere özel kaynakları dinamik hale getirmek ve Türkçe/küresel premium haber yayınlarını entegre etmek için tasarlanan yeni haber çekme mimarisini tanımlar.

---

## 1. Hedef ve Kapsam

### Sorun Tanımı:
* Haber motoru statik olarak sadece HackerNews ve NYT Technology kaynaklarından beslenmektedir.
* Her kategori için katı İngilizce anahtar kelime eşleştirmesi yapıldığından, yerel (Türkçe) haberler veya genel donanım/mobil haberleri sisteme düşmemekte, kategoriler boş kalmaktadır.
* Veritabanında her kategorinin kendine ait olan `rss_feeds` alanı tamamen göz ardı edilmektedir.

### Çözüm Hedefleri:
* Haber çekme motorunun her kategorinin veritabanında kayıtlı olan özel RSS beslemelerini taranması sağlanacaktır.
* Veritabanındaki kategoriler, yüksek kaliteli Türkçe ve global premium RSS kaynakları ile tohumlanacaktır (seeding).
* **Kaynak Bazlı Akıllı Filtreleme (Smart Keyword Filtering Bypass):** 
  * Eğer bir RSS beslemesi kategoriye özel ise (örn: Webrazzi -> Türkiye Haberleri), anahtar kelime süzgeci tamamen devre dışı bırakılarak tüm haberler doğrudan kategoriye dahil edilecektir.
  * Paylaşımlı genel kaynaklardan (HackerNews gibi) gelen haberler ise kategorilere dağıtılmak üzere anahtar kelime filtresine girmeye devam edecektir.

---

## 2. Mimari Tasarım ve Değişiklikler

### A. Port Arayüz Güncellemesi (`internal/domain/port/content_fetcher.go`)
`ContentFetcher` arayüzünün `FetchByKeywords` metodu, haber çekme motorunun kategorinin tüm verilerine (anahtar kelimeler ve özel RSS kaynakları dahil) tam erişebilmesi adına `Fetch` olarak güncellenecektir:

```go
type ContentFetcher interface {
	Fetch(ctx context.Context, topic *entity.Topic) ([]*entity.Article, error)
	SourceType() valueobject.SourceType
}
```

### B. Akıllı RSS Arama Motoru (`internal/adapter/fetcher/rss.go`)
* RSS fetcher içindeki statik `feeds` dizisi, yedekleme (fallback) amacıyla kullanılacaktır.
* `Fetch` metodu içinde, paylaşımlı genel kaynaklar listesi tanımlanacaktır:
  ```go
  var generalFeeds = map[string]bool{
      "https://news.ycombinator.com/rss":                            true,
      "https://hnrss.org/frontpage":                                 true,
      "https://dev.to/feed":                                         true,
      "https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml": true,
  }
  ```
* Her bir besleme taranırken, beslemenin genel bir havuz olup olmadığı (`isShared := generalFeeds[feedURL]`) tespit edilecek ve `parseFeed` metoduna parametre olarak geçilecektir.
* `parseFeed` metodu, eğer `applyFilter` değeri `false` ise anahtar kelime süzgecini es geçerek haberleri doğrudan toplayacaktır.

### C. HackerNews Arama Motoru (`internal/adapter/fetcher/hackernews.go`)
* `Fetch` metodu implemente edilecek ve gelen `topic.Keywords` listesini kullanarak HackerNews API'si üzerinden arama yapmaya devam edecektir.

### D. İş Mantığı Katmanı (`internal/usecase/fetch_articles.go`)
* `Execute` iş akışı basitleştirilecek ve fetcher'ların doğrudan `Fetch(ctx, topic)` metodu tetiklenecektir.

---

## 3. Premium RSS Kaynak Tohumlaması (`backend/cmd/seed_all/main.go`)

Veritabanındaki kategoriler aşağıdaki canlı ve prestijli RSS beslemeleri ile güncellenecektir:

1. **Yapay Zeka (AI):**
   * `https://www.technologyreview.com/topic/artificial-intelligence/feed/` (MIT Technology Review AI)
   * `https://www.wired.com/feed/tag/ai/latest/tgx` (Wired AI)
   * `https://www.artificialintelligence-news.com/feed/` (AI News)
   * *Genel: HackerNews (Kelime Filtreli)*

2. **Mobil:**
   * `https://www.phonearena.com/feed` (PhoneArena)
   * `https://9to5mac.com/feed/` (9to5Mac - Apple)
   * `https://9to5google.com/feed/` (9to5Google - Android)

3. **Uygulama/Yazılım:**
   * `https://dev.to/feed` (Dev.to)
   * `https://feed.infoq.com/` (InfoQ)
   * `https://www.smashingmagazine.com/feed/` (Smashing Magazine)

4. **Donanım:**
   * `https://www.tomshardware.com/feeds/all` (Tom's Hardware)
   * `https://www.anandtech.com/rss/` (AnandTech)

5. **Global Haberler:**
   * `https://techcrunch.com/feed/` (TechCrunch)
   * `https://www.theverge.com/rss/index.xml` (The Verge)
   * `https://www.wired.com/feed/rss` (Wired)

6. **Türkiye Haberleri:**
   * `https://webrazzi.com/feed/` (Webrazzi)
   * `https://shiftdelete.net/feed` (ShiftDelete.net)
   * `https://www.webtekno.com/rss.xml` (Webtekno)
   * `https://www.donanimhaber.com/rss/tumhaberler.xml` (DonanımHaber)

---

## 4. Doğrulama ve Test Planı

1. **Birim Testleri:** `internal/adapter/handler/admin_handler_test.go` ve diğer mock testlerindeki `ContentFetcher` arayüz tanımları yeni `Fetch` imzasına uygun olarak güncellenecek, `go test ./...` komutuyla testlerin tamamı doğrulanacaktır.
2. **Manuel Doğrulama:** 
   * Veritabanı seeding aracı (`go run cmd/seed_all/main.go`) çalıştırılarak tüm yeni kaynaklar yüklenecektir.
   * Admin panelindeki **"Hemen Haber Çek & İşle"** butonuna basılarak Türkiye Haberleri dahil tüm kategorilere pürüzsüzce haberlerin çekildiği gözlemlenecektir.
