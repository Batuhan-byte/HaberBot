---
trigger: on_change
glob: frontend/src/**/*.{js,jsx,css}
description: HaberBot Sayfa Düzeni, Navigasyon ve Drawer Kuralları. Full-width akışkan navbar yerleşimi ve z-index stacking context hatalarını engelleyen fixed drawer konumlandırma kuralları.
---

# 📐 HaberBot Yerleşim, Navigasyon & Drawer Güvenliği (frontend-layout-nav.md)

Bu dosya, HaberBot arayüzünün genel yerleşim (layout) yapısını, responsive ızgaralarını ve kayan yan menü (sliding drawer), karartma (overlay) ve modal gibi `position: fixed` elemanların CSS yerleşim kurallarını tanımlar. Yerleşim şablonları oluştururken veya navigasyon elemanlarını kodlarken bu kurallara kesinlikle uymalısınız.

---

## 🧭 1. Full-Width Fluid (Akışkan) Navbar Standardı

HaberBot'un marka kimliği ve kullanıcı deneyimi gereği, üst navigasyon çubuğu (navbar) **her zaman ekran genişliğinin tamamını kaplamalıdır**.

* **Kural:** Navbar ana konteynerine kesinlikle `max-width: 1200px` veya `margin: 0 auto` gibi genişliği daraltıcı kurallar uygulamayın.
* **Akışkanlık:** Marka logosunu ve kontrol butonlarını ekranın en uç kenarlarına kadar itmek için navbar'ı fluid tasarlayın ve iç kenar boşluğu (padding) verin.

### 🛠️ Fluid Navbar CSS Şablonu:
```css
.webtekno-nav {
  width: 100%;
  max-width: 100%;
  height: 64px;
  background: rgba(5, 5, 5, 0.8);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--border-subtle);
  padding: 0 24px; /* Logoyu sola, butonları en sağa iter */
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 100;
}
```

---

## 🛡️ 2. Stacking Context (Katmanlanma) & Fixed Drawer Yerleşim Kuralları

CSS'te `backdrop-filter`, `transform`, `filter`, `opacity` veya `perspective` gibi özellikler uygulanan konteynerler **yeni bir stacking context (katmanlanma bağlamı) oluşturur**. 

### ⚠️ Stacking Context Hatası Nedir?
Eğer `position: fixed` olan bir tam ekran yan menüyü (`.webtekno-left-drawer`) veya bir arka plan karartmasını (`.drawer-overlay`) bu tür bir ebeveyn div'in (örneğin `<header>` veya `.webtekno-nav`) içine yazarsanız, tarayıcı bu fixed elemanın boyutunu ve konumunu tüm ekran yerine **sadece o ebeveyn elementin boyutuyla (örn: 64px yükseklik) sınırlar**. Menü ekranın tamamına yayılamaz veya z-index değerleri bozulur.

### 🔴 ÇÖZÜM KURALI (Fixed Elemanları En Dışa Koyun):
Kayan yan menüler (drawers), karartma katmanları (overlays) veya modallar **asla `<header>` veya navbar div'inin içerisine gömülmemelidir**. Bunlar her zaman DOM ağacının en dış katmanında, sayfa bileşeninin en üst seviyesinde kardeş eleman (sibling) olarak yer almalıdır.

### 🛠️ Doğru JSX Yapısı Şablonu:
```jsx
import React, { useState } from 'react';
import { Header } from './Header';
import { SidebarDrawer } from './SidebarDrawer';

export function Layout({ children }) {
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);

  return (
    <>
      {/* 1. Header: Sadece kendi sınırları (64px) içinde kalır */}
      <Header onMenuClick={() => setIsDrawerOpen(true)} />

      {/* 2. Drawer ve Overlay: Header'ın DIŞINDA, en üst seviyede kardeş eleman olarak konumlanır */}
      {isDrawerOpen && (
        <>
          {/* Ekran karartması */}
          <div 
            className="drawer-overlay" 
            onClick={() => setIsDrawerOpen(false)} 
          />
          {/* Sol menü */}
          <aside className="webtekno-left-drawer">
            <div className="drawer-header">
              <h3>HaberBot Menü</h3>
              <button onClick={() => setIsDrawerOpen(false)}>Kapat</button>
            </div>
            <nav className="drawer-nav">
              {/* Menü Linkleri */}
            </nav>
          </aside>
        </>
      )}

      {/* 3. Ana İçerik */}
      <main className="main-content">
        {children}
      </main>
    </>
  );
}
```

---

## 📐 3. Responsive Grid (Izgara) Tasarımı

Haber kartlarımızın farklı cihazlarda kırılmadan, akıcı bir şekilde yerleşmesini sağlamak için CSS Grid auto-fit kalıbını kullanın:

### 🛠️ Grid CSS Şablonu:
```css
.articles-grid {
  display: grid;
  /* Ekran genişliğine göre kart sayısını otomatik belirler. Kartlar en az 280px, en fazla 1fr olur. */
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
  width: 100%;
  padding: 24px 0;
}

@media (max-width: 640px) {
  .articles-grid {
    grid-template-columns: 1fr; /* Mobil ekranlarda tek sütuna düşürür */
    gap: 16px;
  }
}
```

---

## 🚫 4. Yerleşim Anti-Patterns (Yapılmaması Gerekenler)

* **Layout Elemanlarına Sabit Piksel Yükseklikleri Vermeyin:** İçeriklerin taşmasını veya kaymasını engellemek için `main-content` veya kart gövdelerine asla sabit yükseklik (`height: 500px`) tanımlamayın. Yükseklik her zaman esnek olmalı, gerekirse iç kenar boşluklarıyla (padding) veya `min-height` ile kontrol edilmelidir.
* **Mobil Uyumluluğu Sadece Medya Sorgularına Bırakmayın:** Medya sorgularını yazmadan önce flex-wrap, CSS grid auto-fit ve yüzdesel genişlikleri kullanarak yerleşimi zaten esnek ve akışkan (fluid) kurgulayın. Medya sorgularını sadece son dokunuşlar için kullanın.
