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

---

## Kullanıcı Profil & Etkileşim Uç Noktaları

### Kamuoyuna Açık Kullanıcı Profili Getir
Bir kullanıcının profil detaylarını, onaylanmış makale istatistiklerini ve (eğer istek atan kişi giriş yapmışsa) sık kullanılanlara eklenme durumunu döndürür.

**GET** `/users/:username`

**Header:**
- `Authorization` (isteğe bağlı): `Bearer <access_token>`

**Yanıt (200 OK):**
```json
{
  "id": "user-uuid",
  "username": "john_doe",
  "email": "john@example.com",
  "role": "User",
  "bio": "Teknoloji meraklısı ve geliştirici.",
  "avatar_url": "/public/avatars/user-user-uuid.png",
  "join_date": "2026-05-29T12:00:00Z",
  "articles_published": 12,
  "is_favorited": true
}
```

### Kendi Profilini Güncelle
Giriş yapmış kullanıcının biyografi metnini veya avatar görselini (dosya veya hazır preset) günceller.

**PUT** `/users/profile`  *(Multipart Form Data)*

**Header:**
- `Authorization`: `Bearer <access_token>`

**Parametreler (Form-Data):**
- `bio` (isteğe bağlı, maks 160 karakter): Yeni biyografi yazısı. HTML karakterleri otomatik kaçırılır.
- `avatar` (isteğe bağlı, dosya): JPG, PNG veya WebP biçimli maks 2MB boyutunda görsel dosyası.
- `preset_avatar` (isteğe bağlı, metin): Hazır avatar seçeneği (`preset-1` ile `preset-8` arası).

**Yanıt (200 OK):**
```json
{
  "id": "user-uuid",
  "username": "john_doe",
  "email": "john@example.com",
  "bio": "Güncellenmiş biyografi metni.",
  "avatar_url": "/public/avatars/user-user-uuid.png",
  "join_date": "2026-05-29T12:00:00Z",
  "created_at": "...",
  "updated_at": "..."
}
```

### Sık Kullanılanlara Ekle
Belirtilen kullanıcıyı giriş yapmış kullanıcının favoriler listesine ekler.

**POST** `/users/:userId/favorites`

**Header:**
- `Authorization`: `Bearer <access_token>`

**Yanıt (201 Created):**
```json
{
  "message": "Kullanıcı favorilerinize eklendi.",
  "favorite_user_id": "target-user-uuid"
}
```

### Sık Kullanılanlardan Çıkar
Belirtilen kullanıcıyı giriş yapmış kullanıcının favorilerinden siler.

**DELETE** `/users/:userId/favorites`

**Header:**
- `Authorization`: `Bearer <access_token>`

**Yanıt (204 No Content):** (İçerik yok)

### Sık Kullanılanlar Listemi Getir
Giriş yapmış kullanıcının favori listesindeki kullanıcıların sayfalanmış bir listesini döndürür.

**GET** `/users/me/favorites`

**Header:**
- `Authorization`: `Bearer <access_token>`

**Sorgu Parametreleri (Query Parameters):**
- `limit` (isteğe bağlı, varsayılan 20): Döndürülecek kullanıcı sayısı.
- `offset` (isteğe bağlı, varsayılan 0): Atlanacak kayıt sayısı.

**Yanıt (200 OK):**
```json
{
  "favorites": [
    {
      "id": "target-user-uuid",
      "username": "target_user",
      "email": "target@example.com",
      "role": "User",
      "bio": "...",
      "avatar_url": "..."
    }
  ],
  "total": 1
}
```

### Profil Şikayet Et
Bir kullanıcının profili hakkında şikayette bulunur. 24 saatlik rate-limiting sınırı mevcuttur.

**POST** `/profile-reports`

**Header:**
- `Authorization`: `Bearer <access_token>`

**İstek (Body):**
```json
{
  "reported_user_id": "target-user-uuid",
  "reason": "spam", // Seçenekler: 'spam', 'offensive_content', 'harassment', 'other'
  "comment": "Şikayet açıklaması..."
}
```

**Yanıt (210 Created):**
```json
{
  "message": "Şikayetiniz başarıyla iletildi. Moderasyon ekibi inceleyecektir."
}
```

---

## Yorum Uç Noktaları

### Yorumları Listele
Belirli bir makaleye ait yorumları çeker.

**GET** `/comments`

**Sorgu Parametreleri:**
- `article_id` (gerekli): Yorumları listelenecek makalenin benzersiz kimliği.

**Yanıt (200 OK):**
```json
{
  "comments": [
    {
      "id": "comment-uuid",
      "article_id": "article-uuid",
      "user_id": "user-uuid",
      "content": "Çok başarılı bir içerik, teşekkürler!",
      "created_at": "2026-05-29T13:00:00Z"
    }
  ]
}
```

### Yorum Gönder
Bir makaleye yeni bir yorum yazar. Kullanıcının platform aktivite puanını ve stats önbelleğini tetikler.

**POST** `/comments`

**Header:**
- `Authorization`: `Bearer <access_token>`

**İstek (Body):**
```json
{
  "article_id": "article-uuid",
  "content": "Harika bir gelişme."
}
```

**Yanıt (201 Created):**
```json
{
  "id": "comment-uuid",
  "article_id": "article-uuid",
  "user_id": "user-uuid",
  "content": "Harika bir gelişme.",
  "created_at": "..."
}
```

