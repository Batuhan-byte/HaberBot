---
trigger: on_change
glob: frontend/src/**/*.{js,jsx,css}
description: HaberBot Zengin HTML İçerik Render ve CSS Formatlama Kuralları. RSS kaynaklarından gelen ham HTML verilerinin güvenli render edilmesi ve .news-content editorial CSS şablonları.
---

# ✍️ Zengin HTML İçerik ve Makale Görünümü (frontend-content.md)

HaberBot, farklı RSS kaynaklarından (HackerNews, Webtekno vb.) gelen ve içerisinde ham HTML etiketleri barındıran haber metinlerini işler. Bu metinlerin okuyucuya **editorial (dergi/gazete kalitesinde)**, son derece şık ve temiz bir tipografiyle sunulması gerekir. Bu kural dosyası, zengin içeriklerin render edilmesi ve CSS formatlama kurallarını tanımlar.

---

## 🔒 1. Güvenli HTML Yorumlama ve Render Şablonu

Backend'den gelen ham HTML payloads (makale içerikleri), React'in korumalı yapısı nedeniyle doğrudan düz metin olarak görünür. Bunları React arayüzünde yorumlayarak göstermek için **React'in dahili interpreter'ını** şık bir sarmalayıcı div ile kullanmalısınız.

### 🛠️ JSX Render Şablonu:
```jsx
import React from 'react';
import './ArticleContent.css'; // .news-content CSS stillerini barındırır

export function ArticleContent({ article }) {
  // Orijinal içerik yoksa fallback olarak özet veya boşluk atanır.
  const rawHtmlContent = article.original_content || article.content || 'İçerik bulunmamaktadır.';

  return (
    <div className="article-content-wrapper">
      <div 
        className="news-content"
        dangerouslySetInnerHTML={{ __html: rawHtmlContent }}
      />
    </div>
  );
}
```

---

## 🎨 2. Zengin CSS Formatlayıcı (.news-content) Standartları

Bileşen içine enjekte edilen HTML etiketlerinin (paragraflar, başlıklar, listeler, alıntılar vb.) HaberBot'un premium koyu temasına kusursuz uyum sağlaması için `index.css` veya ilgili CSS dosyasında **`.news-content`** seçicisi altında şu kurallar tanımlanmalı ve uygulanmalıdır:

### 🛠️ Editorial CSS Format Şablonu:
```css
/* Ana Metin Konteyneri */
.news-content {
  font-family: 'Inter', -apple-system, sans-serif;
  font-size: 17px;
  line-height: 1.85; /* Okuma akışını kolaylaştırmak için yüksek satır aralığı */
  color: #e5e7eb;    /* Yumuşak kırık beyaz, gözü yormaz */
}

/* Paragraflar */
.news-content p {
  margin-bottom: 24px; /* Paragraflar arasında net dikey boşluklar */
  letter-spacing: -0.003em;
}

/* Alt Başlıklar */
.news-content h2, 
.news-content h3 {
  font-family: 'Outfit', sans-serif;
  font-weight: 700;
  color: #ffffff;
  margin-top: 40px;
  margin-bottom: 16px;
  line-height: 1.4;
}

.news-content h2 { font-size: 24px; }
.news-content h3 { font-size: 20px; }

/* Linkler */
.news-content a {
  color: var(--accent-glow);
  text-decoration: none;
  border-bottom: 1px dashed var(--accent-glow);
  transition: all 0.2s ease;
}

.news-content a:hover {
  color: #ffffff;
  border-bottom-style: solid;
}

/* Blok Alıntılar (Blockquotes) */
.news-content blockquote {
  margin: 32px 0;
  padding: 16px 24px;
  background: rgba(255, 255, 255, 0.02);
  border-left: 4px solid var(--brand-blue); /* Sol kenarda mavi vurgu çizgisi */
  border-radius: 0 8px 8px 0;
  font-style: italic;
  color: #d1d5db;
}

/* Kod Blokları (Code & Pre) */
.news-content pre {
  background: #09090b;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.05);
  overflow-x: auto; /* Kodun taşmasını önleyen yatay scroll */
  margin: 24px 0;
}

.news-content code {
  font-family: 'Fira Code', 'Courier New', monospace;
  font-size: 14px;
  color: #3b82f6;
  background: rgba(59, 130, 246, 0.08);
  padding: 2px 6px;
  border-radius: 4px;
}

.news-content pre code {
  background: transparent;
  padding: 0;
  color: #e4e4e7;
}

/* Listeler (Lists) */
.news-content ul, 
.news-content ol {
  margin-bottom: 24px;
  padding-left: 24px;
}

.news-content li {
  margin-bottom: 12px;
}

/* Resimler (Images) */
.news-content img {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 32px 0;
  border: 1px solid rgba(255, 255, 255, 0.05);
}
```

---

## 🚫 3. İçerik Görünümü Anti-Patterns (Yapılmaması Gerekenler)

* **Asla `style` Özniteliklerini Olduğu Gibi Bırakmayın:** Dış kaynaklardan (RSS) kopyalanan HTML metinlerinin içinde yerleşik `style="color: black; font-size: 12px;"` gibi nitelikler olabilir. Bunların HaberBot koyu temasını bozmasını engellemek için, arayüzde gerekirse bu yerleşik stilleri CSS'teki `!important` kurallarıyla veya backend sanitizer seviyesinde temizleyerek ezin.
* **Resim Taşmalarına Dikkat Edin:** `.news-content img` etiketinde `max-width: 100%` tanımının eksik olması, büyük görsellerin makale panelinden dışarı taşarak mobil düzeni bozmasına sebep olur. Bu kuralı kesinlikle atlamayın.
