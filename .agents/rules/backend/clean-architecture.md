---
trigger: on_change
glob: backend/internal/**/*.go
description: HaberBot Clean Architecture ve Katman Sınır Kuralları. Domain katmanı bağımsızlığı, UseCase statelessness ve somut adapter bağımlılık kısıtlamaları.
---

# 🧱 Clean Architecture & Katman Sınırları (rules/backend/clean-architecture.md)

Bu kural dosyası, HaberBot backend katmanının sürdürülebilirliğini, yüksek test edilebilirliğini ve dış kütüphanelerden/veri tabanlarından bağımsızlığını (decoupling) korumak amacıyla tasarlanmıştır. `internal/` altındaki herhangi bir Go paketinde değişiklik yaparken bu kurallara kesinlikle uymalısınız.

*Bu dosya `architecture-patterns` ve `api-and-interface-design` yetenekleri standartlarına göre yapılandırılmıştır.*

---

## 🛑 1. Domain Katmanı Bağımsızlığı (internal/domain/)

Domain katmanı uygulamanın kalbidir ve tamamen saf Go yapılarından oluşmalıdır.

* **SIFIR Bağımlılık Kuralı:** `internal/domain/` altındaki hiçbir dosya, projenin diğer katmanlarından (UseCase, Adapter, Infrastructure) veya dış web/veri tabanı framework'lerinden (Fiber, pgx, gorm vb.) **kesinlikle hiçbir paket ithal edemez (import)**.
* **Ports (Arayüz sözleşmeleri):** Tüm veritabanı yazma/okuma veya dış servis API çağrıları için arayüzler (Go interfaces) `internal/domain/port/` altında tanımlanır.

### 🛠️ Doğru Port Arayüzü Tanımlama Şablonu:
```go
// internal/domain/port/article_repository.go
package port

import (
	"context"
	"HaberBot/internal/domain/entity"
)

// ArticleRepository, makale veritabanı işlemlerinin port sözleşmesidir.
type ArticleRepository interface {
	Save(ctx context.Context, article *entity.Article) error
	GetByID(ctx context.Context, id int64) (*entity.Article, error)
	ListLatest(ctx context.Context, limit int) ([]*entity.Article, error)
}
```

---

## 🧠 2. Durumsuz UseCase Standartları (internal/usecase/)

UseCase katmanı iş mantığını koordine eder (orchestrate) ve tamamen **durumsuz (stateless)** olmalıdır.

* **Statelessness (Durumsuzluk) Kuralı:** UseCase struct'ları içinde asla anlık istek verisi (request-scoped data), geçici veritabanı kayıtları veya transaction durumları gibi state durumları **tutulamaz (struct field olarak saklanamaz)**.
* **Sadece Bağımlılıklar:** UseCase alanları sadece port arayüzleri gibi bağımlılıkları ve gerekirse statik konfigürasyonları barındırabilir. İstek verileri UseCase metodlarına parametre olarak geçilmelidir.

### 🛠️ Doğru UseCase Yapılandırma Şablonu:
```go
// internal/usecase/summarize_article.go
package usecase

import (
	"context"
	"fmt"
	"HaberBot/internal/domain/port"
)

type SummarizeUseCase struct {
	articleRepo port.ArticleRepository // Sadece port bağımlılığı (arayüz)
	aiProcessor port.AIProcessor       // Sadece port bağımlılığı (arayüz)
}

func NewSummarizeUseCase(r port.ArticleRepository, ai port.AIProcessor) *SummarizeUseCase {
	return &SummarizeUseCase{
		articleRepo: r,
		aiProcessor: ai,
	}
}

// Execute metodu durumsuzdur, tüm veri parametre olarak akar.
func (u *SummarizeUseCase) Execute(ctx context.Context, articleID int64) (string, error) {
	article, err := u.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return "", fmt.Errorf("makale bulunamadı: %w", err)
	}

	summary, err := u.aiProcessor.Summarize(ctx, article.OriginalContent)
	if err != nil {
		return "", fmt.Errorf("özet oluşturulamadı: %w", err)
	}

	return summary, nil
}
```

---

## 🚫 3. Somut Sınıf İthalat Yasağı (Concrete Import Restriction)

İş mantığı ve altyapı implementasyonları arasındaki bağları gevşek (decoupled) tutmak amacıyla:

* **Kısıtlama:** `internal/usecase/` ve `internal/adapter/handler/` paketleri, `internal/adapter/repository/` veya `internal/adapter/gateway/` altındaki somut tipleri (struct'ları) **ithal edemez**.
* **Çözüm:** Somut tiplere doğrudan erişim yerine her zaman `port` arayüzü referans alınmalıdır.
```go
// ❌ YANLIŞ (Somut veritabanı implementasyonuna doğrudan bağımlılık)
import "HaberBot/internal/adapter/repository"
func NewUseCase(r *repository.PostgresArticleRepo) ...

//  DOĞRU (Domain port arayüzüne bağımlılık)
import "HaberBot/internal/domain/port"
func NewUseCase(r port.ArticleRepository) ...
```

---

## 🚫 4. Mimari Anti-Patterns (Yapılmaması Gerekenler)

* **UseCase İçinde Fiber Context (`c *fiber.Ctx`) Kullanmayın:** UseCase katmanı Fiber, Gin veya standart `http.Request` gibi HTTP yönlendirici teknolojilerinden tamamen bağımsız olmalıdır. Parametre olarak her zaman standart Go `context.Context` kullanılmalıdır.
* **Domain Entity'lerinde JSON/Veritabanı Tagleri Bulundurmayın:** `internal/domain/entity` içindeki Go struct'ları tamamen temiz olmalı, `gorm:""` veya `json:""` gibi altyapıya bağımlı etiketler içermemelidir. Gerekirse bu etiketler Adapter katmanındaki DTO (Data Transfer Object) struct'larında tanımlanıp entity'ye dönüştürülmelidir (ya da pgx gibi tag gerektirmeyen veri okuyucular tercih edilmelidir).
