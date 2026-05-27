---
trigger: on_change
glob: backend/internal/adapter/gateway/**/*.go
description: HaberBot Yapay Zeka ve Dış API Entegrasyon Kuralları. API kotaları için Uniform Delay (429 engelleme), UseCase B-Planı fallbacks ve prompt mühendisliği güvenlik kuralları.
---

# 🔌 Dış API & Yapay Zeka Entegrasyon Kuralları (rules/backend/api-gateway.md)

Bu kural dosyası, Google Gemini ve Mistral AI gibi dış yapay zeka sağlayıcılarının API sınırlarını (rate limits), kota aşımlarını ve teknik arıza durumlarını en verimli, güvenli ve çökmeye karşı dayanıklı (resilient) şekilde yönetmek amacıyla tasarlanmıştır. AI Gateway adaptörlerinde (`internal/adapter/gateway/`) veya bunları çağıran UseCase'lerde çalışırken bu kurallara kesinlikle uymalısınız.

---

## ⏳ 1. Oran Sınırları ve Tekdüze Bekleme (Uniform Delay for 429 Errors)

Dış yapay zeka API'lerinin ücretsiz veya düşük ücretli katmanları genellikle RPM (Dakika Başına İstek) kotalarına tabidir (örneğin Gemini ücretsiz katmanında 5 RPM sınırı vardır). Bu kotaları aşmamak ve `429 Too Many Requests` hatasını engellemek için şu kurallar geçerlidir:

* **🔴 TEKDÜZE BEKLEME (Uniform Delay) KURALI:**
  Yapay zeka API'lerini bir döngü veya kuyruk içinde çağırırken, **hem başarılı sonuç yollarında hem de hata/başarısızlık durumlarında (failure paths)** döngü başa dönmeden önce mutlaka bekleme (örn: `time.Sleep(12 * time.Second)`) uygulanmalıdır.
  * *Neden?* Hata durumunda beklemeden bir sonraki döngü adımına geçilmesi, saniyeler içinde ardışık hatalar üreterek API kotasının tamamen kilitlenmesine ve cascading 429 hatalarının tüm sistemi kilitlemesine yol açar.

### 🛠️ Doğru Döngü Çağrı Şablonu:
```go
// internal/usecase/daily_pipeline.go
package usecase

import (
	"context"
	"log"
	"time"
)

func (u *DailyPipelineUseCase) ProcessPendingArticles(ctx context.Context) {
	articles, _ := u.repo.GetPending(ctx)

	for _, article := range articles {
		// İşlem adımları...
		err := u.processSingleArticle(ctx, article)
		if err != nil {
			log.Printf("Hata: Makale işlenemedi (%d): %v", article.ID, err)
			// 🔴 HATA DURUMUNDA DA BEKLE: API'yi yormamak için uyuyoruz
			time.Sleep(12 * time.Second)
			continue
		}

		log.Printf("Başarı: Makale işlendi (%d)", article.ID)
		// 🟢 BAŞARILI DURUMDA DA BEKLE
		time.Sleep(12 * time.Second)
	}
}
```

---

## 🛡️ 2. Zarif B-Planı Fallback Mekanizmaları (Graceful Fallback)

Yapay zeka servisinin özetleme veya çeviri yapamaması, tüm makale çekme akışını (pipeline) durdurmamalı ve veritabanı işlemlerini çökertmemelidir.

* **Default Değer Atama Kuralı:** Eğer yapay zeka API'si bir makale için Türkçe özet oluşturmakta başarısız olursa (örneğin kota aşımı veya format hatası), işlemi yarıda kesmek yerine makaleyi yine de veritabanına kaydedin ve özet alanına kullanıcı dostu varsayılan bir yedek dize atayın:
  ```go
  // gateway.Translate/Summarize hata döndürdüğünde UseCase içinde uygulanacak fallback:
  var summaryTR string
  summaryTR, err = u.aiProcessor.Summarize(ctx, article.Content)
  if err != nil {
      // Hata olsa dahi işlemi iptal etmiyoruz, varsayılan bir Türkçe açıklama atıyoruz
      summaryTR = "İçerik teknik olarak yetersiz veya servis yoğunluğu nedeniyle şu an Türkçe özet üretilemedi."
  }
  ```
* **Kırık Kayıtları Gizlememe:** API hataları yüzünden makalelerin veritabanına kaydedilmesini engellemeyin. Arayüzün dynamic property fallbacks (B-Planı) kuralları sayesinde, eksik özetli makaleler dahi arayüzde kırılmadan şık bir şekilde render edilecektir.

---

## ✍️ 3. İstem Güvenliği ve Edge-Cases Prompt Tasarımı

Yapay zeka modellerinden tutarlı ve JSON/metin formatında doğru çıktı alabilmek için sistem istemlerinin (system prompts) uç durumları açıkça tanımlaması gerekir.

* **Prompt Edge-Cases Kuralları:** Yazdığınız prompt'larda şu 3 durum için yapay zekaya ne yapması gerektiğini açıkça dikte edin:
  1. **İçerik çok kısaysa veya sadece tek bir cümleyse:** Ayrıntılı özet üretmek yerine orijinal içeriği doğrudan geri döndürmesini söyleyin.
  2. **İçerik zaten Türkçe ise:** Yeniden çeviri yapmak yerine metni doğrudan korumasını belirtin.
  3. **İçerik teknik olarak anlamsız veya kod bloğundan ibaretse:** Hata fırlatmak yerine "İçerik teknik dokümantasyon içerdiğinden özetlenemedi" gibi kararlı ve standart bir çıktı vermesini dikte edin.

---

## 🚫 4. API Entegrasyon Anti-Patterns (Yapılmaması Gerekenler)

* **Asla API Anahtarlarını Koda Gömerek Kaydetmeyin (No Hardcoded API Keys):** Dış API anahtarları (`AI_API_KEY`) kesinlikle Go kodlarının içine doğrudan yazılmamalıdır. Her zaman `.env` dosyasından okunmalı ve `container.go` aracılığıyla enjekte edilmelidir.
* **Gereksiz Arama/Sorgu (Polling) Yapmayın:** Yapay zeka servislerinin kotalarını tüketmemek için, her HTTP isteğinde (örneğin ana sayfa her açıldığında) Gemini veya Mistral API'sini doğrudan çağırmayın. Yapay zeka işlemlerini her zaman **sadece arka plan iş hatlarında (fetch pipeline)** çalıştırıp sonuçları veritabanına yazın; ana arayüz sadece veritabanını okumalıdır.
