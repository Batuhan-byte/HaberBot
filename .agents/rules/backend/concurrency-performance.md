---
trigger: on_change
glob: backend/**/*.go
description: HaberBot Go Eşzamanlılık, Arka Plan Görevleri ve Performans Kuralları. Goroutine sızıntısı engelleme, race conditions koruması ve context iptal yönetimleri.
---

# ⚡ Eşzamanlılık & Performans Kuralları (rules/backend/concurrency-performance.md)

Go dilinin en güçlü yanlarından biri olan eşzamanlılık (concurrency - Goroutines & Channels), hatalı kullanıldığında bellek sızıntılarına (goroutine leaks), yarış durumlarına (race conditions) veya sistem kilitlenmelerine yol açabilir. HaberBot'un veri çekme (cron scheduler) ve paralel yapay zeka işleme hatlarında eşzamanlı kod yazarken bu kurallara kesinlikle uymalısınız.

*Bu dosya `golang-pro` ve `performance-optimization` yetenekleri standartlarına göre yapılandırılmıştır.*

---

## 🛑 1. Goroutine Sızıntılarını Önleme (No Goroutine Leaks)

Bir goroutine başlatıldığında, yaşam döngüsünü tamamlayıp sonlanması garanti edilmelidir. Sonsuza kadar açık kalan goroutine'ler bellek sızıntısına neden olur ve bir süre sonra sunucuyu çökertir.

* **🔴 CONTEXT İLE İPTAL KURALI:**
  Uzun süren veya arka planda koşan tüm eşzamanlı işlemlere parametre olarak **`context.Context`** geçilmelidir. Goroutine içindeki ana işlem döngüleri, her adımda veya belirli aralıklarla `ctx.Done()` kanalını dinleyerek iptal (cancellation) veya zaman aşımı (timeout) sinyali geldiğinde otonom olarak sonlanmalıdır.

### 🛠️ Doğru Arka Plan Goroutine Döngü Şablonu:
```go
// internal/infrastructure/scheduler/cron.go
package scheduler

import (
	"context"
	"log"
	"time"
)

type CronScheduler struct {
	interval time.Duration
}

func (s *CronScheduler) Start(ctx context.Context) {
	// Yeni bir goroutine içinde arka plan işimizi başlatıyoruz
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		log.Println("Arka plan zamanlayıcı başlatıldı...")

		for {
			select {
			// 🔴 KRİTİK ADIM: Ebeveyn context sonlandığında goroutine güvenle kapanır!
			case <-ctx.Done():
				log.Println("İptal sinyali alındı, arka plan zamanlayıcı kapatılıyor.")
				return
			case <-ticker.C:
				// Zamanı gelen arka plan görevini çalıştır
				s.executeTask(ctx)
			}
		}
	}()
}

func (s *CronScheduler) executeTask(ctx context.Context) {
	// İstek/İşlem sınırı için timeout context oluşturulur
	taskCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	log.Println("Zamanlanmış veri çekme görevi başladı...")
	// Görevi taskCtx ile çalıştır...
}
```

---

## 🏁 2. Veri Yarışları (Race Conditions) Engelleme

Birden fazla goroutine'in aynı bellek hücresine (örneğin bir map, slice veya struct field) aynı anda yazmaya veya okumaya çalışması tanımsız davranışlara (undefined behavior) ve çökmelere sebep olur.

* **Güvenli Paylaşım Kuralları:**
  1. **Kanallar (Channels):** "Belleği paylaşarak iletişim kurmayın, iletişim kurarak belleği paylaşın" ilkesine sadık kalın. Eşzamanlı işlemler arasında veri transferi için channels kullanın.
  2. **Karşılıklı Dışlama (sync.Mutex):** Eğer ortak bir değişken (örneğin bir cache map veya sayaç) birden fazla goroutine tarafından okunup yazılıyorsa, okuma/yazma bloklarını mutlaka `sync.RWMutex` veya `sync.Mutex` ile kilitleyin:

### 🛠️ sync.RWMutex ile Güvenli Bellek Paylaşım Şablonu:
```go
package cache

import "sync"

type MemoryCache struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		store: make(map[string]string),
	}
}

func (c *MemoryCache) Set(key, value string) {
	c.mu.Lock()         // Yazma kilidi (Write Lock)
	defer c.mu.Unlock() // İşlem bitiminde kilidi aç
	c.store[key] = value
}

func (c *MemoryCache) Get(key string) (string, bool) {
	c.mu.RLock()         // Okuma kilidi (Read Lock - Birden fazla okuyucu aynı anda erişebilir)
	defer c.mu.RUnlock()
	val, ok := c.store[key]
	return val, ok
}
```

---

## 🚫 3. Eşzamanlılık Anti-Patterns (Yapılmaması Gerekenler)

* **Sonsuz Bloklanan Kanallar Açmayın (No Blocked Channels):** Kanallara veri gönderirken veya alırken buffer kapasitesi belirtilmemişse (unbuffered channels) alıcı veya gönderici hazır olana kadar işlem bloklanır. Bu kilitlenmeleri (deadlocks) önlemek için gerekirse buffered channels kullanın veya `select` içinde `time.After` ile zaman aşımı tanımlayın.
* **Goroutine İçinde Değişken Kapsamlarına Dikkat Edin (Loop Variable Capture):** Go 1.22 öncesi sürümlerde, döngü değişkenlerini goroutine içine doğrudan geçmek yarış hatasına yol açar. Her ihtimale karşı döngü değişkenini goroutine'e parametre olarak geçin veya yerel kopyasını oluşturun:
  ```go
  // ❌ YANLIŞ
  for _, item := range list {
      go func() {
          process(item) // Tüm goroutine'ler aynı bellek adresini okuyabilir (yarış durumu)
      }()
  }

  //  DOĞRU
  for _, item := range list {
      go func(val string) {
          process(val)
      }(item) // Değişken kopyalanarak güvenle aktarılır
  }
  ```
* **İşlem Sınırı Olmadan Goroutine Başlatmayın (No Unbounded Goroutines):** Binlerce RSS kaynağını işlerken, hepsi için aynı anda sınır koymadan goroutine açmak sistem kaynaklarını tüketir. Bu tür büyük toplu işlerde her zaman iş parçacığı havuzu (Worker Pool) veya sınırlayıcı kanallar (Semaphore) kullanarak eşzamanlı çalışan goroutine sayısını sınırlandırın.
