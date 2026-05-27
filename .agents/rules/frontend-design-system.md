---
trigger: on_change
glob: frontend/src/**/*.{js,jsx,css}
description: HaberBot Tasarım Sistemi Kuralları. Renk paleti, kırmızı yasağı, cam efekti (glassmorphism), Outfit/Inter yazı tipleri ve mikro-animasyon CSS kod şablonları.
---

# 🎨 HaberBot Tasarım Sistemi & Estetik Kuralları (frontend-design-system.md)

Bu kural dosyası, HaberBot'un benzersiz, üst düzey ve modern koyu tema estetiğini korumak amacıyla tasarlanmıştır. Görsel veya stil bazlı bir çalışma (CSS, JSX sınıf düzenlemeleri, yeni arayüz bileşenleri oluşturma) yaparken bu dosyadaki şablonları ve kuralları birebir uygulamalısınız.

*Bu dosya `frontend-design` yeteneği standartlarına göre yapılandırılmıştır.*

---

## 🎨 1. Resmi Renk Paleti ve Kırmızı Rengin Kesin Yasağı

HaberBot, koyu temanın asaletini ve mavi ışığın canlılığını yansıtan **mavi ve siyah/gri tonlar** üzerine markalanmıştır. Arayüzün marka kimliğini bozacak jenerik renk geçişlerinden (özellikle yapay zeka tasarımlarında sıkça görülen pembe/mor gradyanlar) ve **kırmızı** renklerden kaçınılmalıdır.

### 🔴 KRİTİK YASAK (Kırmızı Rengin Kesin Yasağı):
HaberBot'un marka kimliği gereği arayüzde **kırmızı ve tonları (kırmızı butonlar, kırmızı borderlar vb.) kesinlikle yasaktır**. Uygulama tamamen mavi renk kimliği altındadır. Sadece ve sadece kritik hata durumlarını gösteren küçük hata etiketlerinde/metinlerinde standart yumuşak bir kırmızı (`rgba(239, 68, 68, 0.9)`) kullanılabilir, onun dışında tüm etkileşimler ve vurgular mavi tonlarında olmalıdır.

### 🎨 Renk Değişkenleri (CSS Variables):
Bileşenleri şekillendirirken `frontend/src/index.css` içindeki resmi CSS değişkenlerini kullanın:

```css
:root {
  --bg-pitch-black: #000000;   /* Derin arka planlar */
  --bg-dark-panel: #050505;    /* Kartlar, yan menüler ve paneller */
  --bg-card-hover: #0d0d0d;    /* Hover sırasındaki panel rengi */
  
  --brand-blue: #2563eb;       /* Resmi HaberBot marka mavisi */
  --accent-glow: #3b82f6;      /* Parlamalar, aktif sınırlar ve vurgular */
  --text-primary: #ffffff;     /* Crisp beyaz okunabilir metinler */
  --text-muted: #a1a1aa;       /* Açıklama ve tarih metinleri (Gümüş Gri) */
  --border-subtle: rgba(255, 255, 255, 0.05); /* Panel ve kart ince sınırları */
}
```

---

## 💎 2. Cam Efekti (Glassmorphism) ve Derinlik

HaberBot panelleri, arka planda hafif bir parıltı ve bulanıklık hissi veren modern cam efektleriyle tasarlanmalıdır.

### 🛠️ Cam Efekti CSS Şablonu:
Yeni bir kart veya panel eklerken aşağıdaki CSS özelliklerini bir arada kullanarak derinlik hissi yaratın:

```css
.premium-glass-card {
  background: rgba(5, 5, 5, 0.75);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--border-subtle);
  border-radius: 12px;
  box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.37);
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
```

---

## ✍️ 3. Premium Tipografi Standartları

Bileşenlerinizde asla tarayıcının varsayılan yazı tiplerini bırakmayın. Font eşleştirmesini her zaman şu şekilde kurun:

* **Display / Başlıklar (Heading Font):** Güçlü, modern ve karakterli bir duruş için **Outfit** (veya destekleyici olarak **Syne**) yazı tipini tercih edin.
* **Gövde Metinleri (Body Font):** Okuma kolaylığı sağlamak ve gözü yormamak için **Inter** veya sistemin temiz **sans-serif** fontlarını tercih edin.

### CSS Kullanım Örneği:
```css
h1, h2, h3, .heading-premium {
  font-family: 'Outfit', 'System-UI', sans-serif;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--text-primary);
}

p, .body-premium {
  font-family: 'Inter', sans-serif;
  font-weight: 400;
  line-height: 1.6;
  color: var(--text-muted);
}
```

---

## 🔄 4. Mikro-Animasyonlar ve Hover Efektleri

Arayüzün canlı ve etkileşime duyarlı olduğunu hissettirmek için tüm tıklanabilir veya üzerine gelinebilir kart/butonlarda mikro-animasyonlar kullanmak **zorunludur**.

### 🛠️ Kart Hover Parlama (Radial Glow) ve Ölçeklendirme (Scale) Şablonu:
Haber kartlarına (ArticleCard) veya menü ögelerine tıklandığında/üzerine gelindiğinde uygulanacak CSS kuralları:

```css
.interactive-card {
  position: relative;
  overflow: hidden;
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1), 
              border-color 0.3s ease, 
              box-shadow 0.3s ease;
}

/* Hover Durumu */
.interactive-card:hover {
  transform: translateY(-4px) scale(1.02); /* Hafif yukarı kayma ve büyüme */
  border-color: rgba(59, 130, 246, 0.3);    /* Vurgulu mavi sınır */
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.5), 
              0 0 24px rgba(59, 130, 246, 0.12); /* Yumuşak mavi dış parlama */
}

/* Tıklama Durumu (Active State Feedback) */
.interactive-card:active {
  transform: translateY(-2px) scale(0.98); /* Tıklama hissi veren geri çekilme */
}
```

---

## 🚫 5. Tasarımsal Anti-Patterns (Yapılmaması Gerekenler)

* **Asla Düz Düz Renkler Kullanmayın:** Tamamen düz siyah (`#000000`) üzerine dümdüz gri sınır çizgileri çekmekten kaçının. Her zaman hafif degrade parlamalar, glassmorphism veya radial degradeler ekleyerek derinlik oluşturun.
* **Yarıçapları (Border-Radius) Karıştırmayın:** Arayüzde bir kart `12px` köşe yuvarlatmasına sahipse, diğer tüm kartlar da `12px` olmalıdır. Butonlar `8px` ise genel olarak tüm küçük butonlar bu standartta kalmalıdır.
* **Anlık Geçişlerden Kaçının (No Instant Jumps):** Arayüzde hiçbir renk veya büyüklük hover durumunda anında değişmemelidir. Her etkileşime en az `transition: all 0.3s cubic-bezier(...)` ekleyin.
