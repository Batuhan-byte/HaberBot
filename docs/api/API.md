# HaberBot REST API Referansı

Temel URL (Base URL): `http://localhost:8080/api/v1`

## Genel (Public) Uç Noktalar

### Son Makaleleri Getir
En son işlenen makalelerin sayfalanmış bir listesini alır.

**GET** `/articles`

**Sorgu Parametreleri (Query Parameters):**
- `limit` (isteğe bağlı, varsayılan 10): Döndürülecek makale sayısı.

**Yanıt (200 OK):**
```json
{
  "articles": [
    {
      "id": "uuid",
      "title": "Original Title",
      "title_tr": "Türkçe Başlık",
      "original_url": "https://...",
      "source": "hackernews",
      "original_content": "...",
      "summary_tr": "Türkçe özet...",
      "score": 100,
      "topic_id": "uuid",
      "image_url": "",
      "processed_at": "2024-01-01T12:00:00Z",
      "fetched_at": "2024-01-01T11:55:00Z",
      "created_at": "2024-01-01T11:55:00Z"
    }
  ]
}
```

### ID'ye Göre Makale Getir
**GET** `/articles/:id`

### Konuları Getir
**GET** `/topics`

### Konuya Göre Makaleleri Getir
**GET** `/topics/:slug/articles`
**Sorgu Parametreleri:**
- `page` (isteğe bağlı, varsayılan 1)
- `limit` (isteğe bağlı, varsayılan 10)

## Yönetici (Admin) Uç Noktaları
*Şu Header'ı Gerektirir: `X-Admin-API-Key: <your_admin_key>`*

### Konu Oluştur
**POST** `/admin/topics`
```json
{
  "name": "Yapay Zeka",
  "slug": "yapay-zeka",
  "keywords": ["ai", "machine learning"],
  "sources": ["hackernews", "rss"],
  "rss_feeds": ["https://..."],
  "is_active": true
}
```

### Konu Güncelle
**PUT** `/admin/topics/:id`

### Konu Sil
**DELETE** `/admin/topics/:id`

### Veri Çekme (Fetch) İşlem Hattını Tetikle
Makale getirme sürecini manuel olarak tetikler.
**POST** `/admin/fetch`

### Yapay Zeka İşleme Hattını Tetikle
Yapay zeka çeviri ve özetleme sürecini manuel olarak tetikler.
**POST** `/admin/process`
