---
trigger: on_change
glob: frontend/src/**/*.{js,jsx}
description: HaberBot Etkileşimli Arayüz Elemanları, Formlar ve Admin Paneli Kuralları. Input focus neon mavi glow efektleri, buton scale-feedback active durum animasyonu, skeleton loader ve empty state şablonları.
---

# 🎯 Etkileşim, Formlar & Admin Kuralları (rules/frontend/interactive-admin.md)

Bu dosya; HaberBot uygulamasındaki kullanıcı etkileşimlerini, form girdi elemanlarını (giriş sayfaları, kaynak ekleme alanları), skeleton loader yükleme ekranlarını, boş/hata durum tasarımlarını ve Admin Dashboard kontrollerini yöneten kuralları ve kod şablonlarını barındırır. Etkileşimli bir UI elemanı tasarlarken bu standartlara uymalısınız.

---

## 📝 1. Form Elemanları ve Input Focus Glow Tasarımı

Kullanıcı form girdileri (Login girdileri, Admin RSS kaynak ekleme alanları) premium koyu temaya tam uyum sağlamalı, focus durumunda şık bir mavi ışık parlaması (neon glow) yaymalıdır.

### 🛠️ Input CSS Şablonu:
```css
.premium-input {
  width: 100%;
  padding: 12px 16px;
  background: rgba(13, 13, 13, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  color: #ffffff;
  font-family: 'Inter', sans-serif;
  font-size: 15px;
  outline: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

/* Focus Durumu (Neon Mavi Parlama) */
.premium-input:focus {
  border-color: var(--accent-glow);
  box-shadow: 0 0 10px rgba(59, 130, 246, 0.35);
}

/* Placeholder Rengi */
.premium-input::placeholder {
  color: rgba(255, 255, 255, 0.3);
}
```

---

## ⚡ 2. Buton Mikro-Etkileşimleri ve Geri Bildirim (Active Scale)

Tıklanabilir tüm butonlarda kullanıcının tıklama aksiyonuna fiziksel geri bildirim hissi verilmelidir:

### 🛠️ Premium Buton CSS Şablonu:
```css
.premium-button {
  background: var(--brand-blue);
  color: #ffffff;
  padding: 12px 24px;
  border-radius: 8px;
  font-family: 'Outfit', sans-serif;
  font-weight: 600;
  border: none;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: background-color 0.2s ease, transform 0.1s ease;
}

/* Hover ve Active Mikro Geri Bildirimi */
.premium-button:hover {
  background: var(--accent-glow);
}

.premium-button:active {
  transform: scale(0.95); /* Tıklama sırasında hafifçe içe çekilerek basılma hissi verir */
}

.premium-button:disabled {
  background: #27272a;
  color: #71717a;
  cursor: not-allowed;
  transform: none;
}
```

---

## ⏳ 3. Premium Skeleton Loader Şablonu

Kullanıcılara sonsuz yükleme animasyonları veya düz "yükleniyor..." yazıları göstermek yerine, HaberBot'un premium kart yapısına uygun parlayıp sönen (pulse efekti) skeleton elemanları gösterilmelidir.

### 🛠️ JSX Skeleton Bileşen Şablonu:
```jsx
import React from 'react';
import './Skeleton.css'; // Pulse animasyon kodunu barındırır

export function ArticleCardSkeleton() {
  return (
    <div className="skeleton-card premium-glass-card">
      <div className="skeleton-line skeleton-image" />
      <div className="skeleton-line skeleton-title" />
      <div className="skeleton-line skeleton-text" />
      <div className="skeleton-line skeleton-text short" />
    </div>
  );
}

export function SkeletonGrid() {
  return (
    <div className="articles-grid">
      {[1, 2, 3, 4, 5, 6].map((i) => (
        <ArticleCardSkeleton key={i} />
      ))}
    </div>
  );
}
```

### 🛠️ Skeleton CSS Şablonu:
```css
.skeleton-card {
  height: 380px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.skeleton-line {
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0.03) 25%,
    rgba(255, 255, 255, 0.08) 50%,
    rgba(255, 255, 255, 0.03) 75%
  );
  background-size: 200% 100%;
  animation: loading-pulse 1.5s infinite ease-in-out;
  border-radius: 4px;
}

.skeleton-image {
  width: 100%;
  height: 180px;
  border-radius: 8px;
}

.skeleton-title {
  width: 85%;
  height: 24px;
}

.skeleton-text {
  width: 100%;
  height: 16px;
}

.skeleton-text.short {
  width: 60%;
}

@keyframes loading-pulse {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
```

---

## 🏜️ 4. Boş (Empty) ve Hata (Error) Durum Tasarımları

Bir kategoride haber bulunmadığında ekranı bomboş bırakmayın. Şık bir ikon ve açıklayıcı buton barındıran bir placeholder gösterin:

### 🛠️ Empty State JSX Şablonu:
```jsx
import React from 'react';

export function EmptyState({ message, onActionClick, actionText }) {
  return (
    <div className="empty-state-container premium-glass-card">
      <div className="empty-icon-glow">⚡</div>
      <h3 className="empty-title">Sonuç Bulunamadı</h3>
      <p className="empty-text">{message || 'Bu kategoriye ait haber henüz eklenmemiştir.'}</p>
      {onActionClick && (
        <button className="premium-button" onClick={onActionClick}>
          {actionText || 'Yeniden Dene'}
        </button>
      )}
    </div>
  );
}
```

```css
.empty-state-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  text-align: center;
  border: 1px solid rgba(255, 255, 255, 0.05);
  margin: 32px 0;
}

.empty-icon-glow {
  font-size: 40px;
  color: var(--accent-glow);
  text-shadow: 0 0 20px rgba(59, 130, 246, 0.4);
  margin-bottom: 16px;
}

.empty-title {
  font-size: 20px;
  margin-bottom: 8px;
}

.empty-text {
  color: var(--text-muted);
  max-width: 320px;
  margin-bottom: 24px;
  font-size: 14px;
}
```

---

## 🛠️ 5. Admin Yönetim Paneli Standartları

Admin Dashboard arayüzünde (`AdminPage`, `AdminContent` vb.) RSS veri çekme tetikleyicileri ve kaynak ayarları gibi gelişmiş ayarlar yer alır.

* **Durum Butonları (Trigger Buttons):** Haber çekme işlemini tetikleyen "Haber Çek (Fetch)" butonları tıklama yapıldığında anında deaktif edilmeli (`disabled`) ve buton içinde loading spinner veya "Çekiliyor..." yazısı belirmelidir. Bu sayede kullanıcının butona art arda tıklayarak sunucu istek hattını kilitlemesi (race condition) engellenir.
* **Tip Güvenliği ve Doğrulama:** Yeni kaynak ekleme formlarında girilen RSS URL alanları regex ile mutlaka doğrulanmalı (`https://` protokolü kontrolü) ve geçersiz URL girildiğinde buton tıklanamaz olmalıdır.
