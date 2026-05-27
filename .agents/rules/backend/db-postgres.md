---
trigger: on_change
glob: backend/internal/adapter/repository/**/*.go
description: HaberBot PostgreSQL ve Veritabanı Repository Kuralları. pgxpool bağlantı havuzu yönetimi, rows.Close() sızıntı engelleme, SELECT * yasağı ve indeks standartları.
---

# 💾 PostgreSQL & Veri Tabanı Kuralları (rules/backend/db-postgres.md)

Bu kural dosyası, PostgreSQL (Neon) veritabanı erişimlerinin yüksek performanslı, optimize edilmiş ve sızıntısız (leak-free) şekilde yapılması amacıyla tasarlanmıştır. Veritabanı repository adaptörlerinde (`internal/adapter/repository/`) herhangi bir SQL sorgusu yazarken veya düzenlerken bu kurallara kesinlikle uymalısınız.

*Bu dosya `neon-postgres` yeteneği standartlarına göre yapılandırılmıştır.*

---

## 🏎️ 1. PostgreSQL Bağlantı Havuzu (pgxpool) Kuralları

HaberBot, yüksek eşzamanlılık gerektiren arka plan veri çekme görevlerine sahip olduğundan, her veritabanı işlemi için yeni bir bağlantı açmak yerine verimli bir bağlantı havuzu (`pgxpool.Pool`) kullanır.

* **Havuz Yapılandırma Standartları:**
  * Maksimum Açık Bağlantı (`MaxConns`): Trafiğe ve sunucu limitlerine uygun olarak **en fazla 25** olmalıdır.
  * Minimum Açık Bağlantı (`MinConns`): Başlangıçta hazır bekleyen bağlantı sayısı **en az 5** olmalıdır.
  * Boşta Kalma Süresi (`MaxConnIdleTime`): Boştaki bağlantılar en fazla 10 dakika (`10m`) sonra kapatılmalıdır.

---

## 🛡️ 2. Bağlantı Sızıntılarını Engelleme (Connection Leaks)

Go ve PostgreSQL entegrasyonlarında en sık karşılaşılan hata, veritabanından veri çektikten sonra `rows` veya `conn` kaynaklarının açık bırakılması ve havuzun kilitlenmesidir.

* **🔴 KRİTİK DEFER rows.Close() KURALI:**
  Her `pool.Query(ctx, ...)` veya `tx.Query(ctx, ...)` çağrısından hemen sonra, hata kontrolünden önce veya hemen sonra **`defer rows.Close()`** çağrılmalıdır. Aksi takdirde, veritabanı bağlantısı havuzu boşalana kadar açık kalır ve sistem bir süre sonra kilitlenir.

### 🛠️ Doğru Repository Sorgu Şablonu:
```go
// internal/adapter/repository/postgres_article.go
package repository

import (
	"context"
	"fmt"
	"HaberBot/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresArticleRepository struct {
	pool *pgxpool.Pool // pgxpool bağlantı havuzu
}

func NewPostgresArticleRepository(pool *pgxpool.Pool) *PostgresArticleRepository {
	return &PostgresArticleRepository{pool: pool}
}

func (r *PostgresArticleRepository) ListLatest(ctx context.Context, limit int) ([]*entity.Article, error) {
	// SELECT * yasağı: kolon isimleri açıkça belirtilir.
	query := `
		SELECT id, title, title_tr, content, summary_tr, created_at 
		FROM articles 
		ORDER BY created_at DESC 
		LIMIT $1`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("sorgu başarısız: %w", err)
	}
	// 🔴 KRİTİK ADIM: rows nesnesi defer ile hemen kapatılır!
	defer rows.Close()

	var articles []*entity.Article
	for rows.Next() {
		var a entity.Article
		err := rows.Scan(&a.ID, &a.Title, &a.TitleTR, &a.Content, &a.SummaryTR, &a.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("satır okuma başarısız: %w", err)
		}
		articles = append(articles, &a)
	}

	// rows döngüsü bittikten sonra oluşabilecek hataların kontrolü
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("satır okuma sırasında hata: %w", err)
	}

	return articles, nil
}
```

---

## 🚫 3. SELECT * Kesin Yasağı (Named Fields)

* **Neden?** Veritabanı şeması genişledikçe veya değiştikçe `SELECT *` kullanımı Go'daki `Scan` işlemini bozar. Ayrıca, ihtiyacımız olmayan büyük HTML veya metin kolonlarını (`content`) her sorguda çekmek gereksiz bellek ve ağ tüketimine neden olur.
* **Kural:** SQL sorgularında **asla `SELECT *` kullanılamaz**. Çekilmek istenen kolon isimleri her zaman SQL cümlesinde açıkça tek tek yazılmalıdır.

---

## 📇 4. İndeks Kapsamı (Index Coverage)

Bileşenlerimizin veya API uçlarımızın filtreleme ve sıralama yaptığı tüm veritabanı kolonları indekslenmiş olmalıdır.

* **Kural:** `WHERE` koşulunda filtreleme yapılan (örneğin `topic_id`) veya `ORDER BY` ile sıralama yapılan (örneğin `created_at`, `fetched_at`) kolonlar veritabanı migrasyonlarında (`migrations/*.sql`) indekslenmiş olmalıdır.
* **Örnek Migrasyon İndeks Tanımı:**
  ```sql
  CREATE INDEX IF NOT EXISTS idx_articles_topic_id ON articles(topic_id);
  CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at DESC);
  ```

---

## 🚫 5. Veritabanı Anti-Patterns (Yapılmaması Gerekenler)

* **Bağlantı Havuzunu Kapatmayı Unutmayın:** Uygulama kapanırken (graceful shutdown) `pool.Close()` çağrılarak açık kalan tüm soketler ve havuz bağlantıları kapatılmalıdır.
* **UseCase veya Handler İçinde Ham SQL Yazmayın:** Hiçbir SQL sorgusu Repository adaptör katmanının dışına taşmamalıdır. İş mantığı katmanı veritabanının yapısından haberdar olmamalıdır.
