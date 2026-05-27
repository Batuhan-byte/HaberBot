---
trigger: always_on
glob: frontend/src/**/*.{js,jsx,css}
description: HaberBot Frontend Kaptan Köşkü (Main Router). React 19 standartları, logic/presentation ayrımı ve alt kural yönlendirme haritası.
---

# ⚓ HaberBot Frontend Kaptan Köşkü (frontend.md)

Bu dosya HaberBot frontend geliştirme süreçlerinin ana kumanda merkezidir. Antigravity frontend ile ilgili bir dosya okuduğunda veya değiştireceğinde **her zaman ilk olarak bu kural dosyasını** referans alır. 

Geliştirme yaparken, üzerinde çalıştığınız konuya ait en detaylı kurallara ve kopyalanabilir şablonlara ulaşmak için aşağıdaki yönlendirme matrisini (routing) kullanmalı ve ilgili alt dosyayı belleğinize yüklemelisiniz.

---

## 🧭 1. Alt Kural Yönlendirme Haritası (Routing Matrix)

Çalıştığınız dosyaya veya konuya göre aşağıdaki bağlantılara giderek o alana özel detaylı kuralları okuyun:

1. **Görsel tasarım, renk paleti, cam efekti, Outfit/Inter yazı tipleri veya hover animasyonları yazarken/düzenlerken:**
   👉 [Tasarım Sistemi & Estetik Kuralları (frontend/design-system.md)](file:///c:/Users/Batuhan/Desktop/deneme%20projem/.agents/rules/frontend/design-system.md)
   
2. **TanStack React Query custom hooks yazarken, veri çekme mantığı kurgularken veya Gemini API kotaları için dayanıklılık (resilience) fallbacks uygularken:**
   👉 [Durum Yönetimi, API & B-Planı Kuralları (frontend/state-api.md)](file:///c:/Users/Batuhan/Desktop/deneme%20projem/.agents/rules/frontend/state-api.md)

3. **Geniş akışkan navbar (Header), responsive ızgaralar (grid) tasarlarken veya fixed sidebar drawer yerleşimleri (stacking context sorunları) yaparken:**
   👉 [Navigasyon & Yerleşim Kuralları (frontend/layout-nav.md)](file:///c:/Users/Batuhan/Desktop/deneme%20projem/.agents/rules/frontend/layout-nav.md)

4. **RSS feed veya makalelerden gelen ham HTML verilerini temiz ve premium bir şekilde (`dangerouslySetInnerHTML`) render ederken veya okuma CSS stilleri yazarken:**
   👉 [Zengin HTML İçerik Okuma Kuralları (frontend/content.md)](file:///c:/Users/Batuhan/Desktop/deneme%20projem/.agents/rules/frontend/content.md)

5. **Form elemanları, login girdileri, premium skeleton loader'lar, boş durumlar veya Admin paneli RSS kaynağı ekleme/çıkarma kontrolleri yazarken:**
   👉 [Etkileşim, Formlar & Admin Kuralları (frontend/interactive-admin.md)](file:///c:/Users/Batuhan/Desktop/deneme%20projem/.agents/rules/frontend/interactive-admin.md)

---

## 🏗️ 2. Core React 19 Geliştirme Standartları

React 19'un yeni getirdiği özellikleri ve davranışları projemizde kusursuz uygulamak için şu kurallara uyun:

### ⚠️ Arayüz Render Hatalarını Engelleme (Ternary vs &&)
React 19'da boolean veya sayısal ifadelerin rendering sırasında beklenmeyen `0` veya boşluk render etmesini engellemek için conditional rendering yaparken **her zaman ternary (`? :`) kullanın**. Asla `&&` kullanmayın.
```javascript
// ❌ YANLIŞ (React 19'da render hatalarına veya '0' çıktısına sebep olabilir)
{articles.length && <ArticleList />}

//  DOĞRU
{articles.length > 0 ? <ArticleList /> : null}
```

### ⚡ useTransition ile Arayüz Blokajlarını Engelleme
Veri yükleme, filtreleme veya arama gibi acil olmayan güncellemelerde kullanıcının etkileşimini (örneğin buton tıklamalarını) engellememek için `useTransition` kullanın:
```javascript
import { useTransition } from 'react';

const [isPending, startTransition] = useTransition();

const handleTabChange = (nextTab) => {
  startTransition(() => {
    setActiveTab(nextTab);
  });
};
```

---

## 🧱 3. Sunum ve Mantık Ayrımı (Presentation & Logic Separation)

HaberBot frontend kodunun okunabilirliğini korumak için **Clean Architecture** prensiplerini frontend'e de yansıtıyoruz:

* **Presentation (Sunum) Katmanı (`components/` ve `pages/`):** 
  * Doğrudan HTTP request (`fetch`, `axios`) yapamaz.
  * Karmaşık veri filtreleme veya caching logic'leri içeremez.
  * Görevi sadece UI render etmek, user input'larını almak ve bu etkileşimleri hook'lara bildirmektir.
* **Logic (Mantık) Katmanı (`hooks/` ve `services/`):**
  * HTTP istekleri `services/api.js` içinde tanımlanır.
  * Bu api istekleri `@tanstack/react-query` ile `hooks/` altındaki custom hook'larda sarmalanır.
  * Bileşenler sadece bu custom hook'ları çağırarak anlık durumları (data, isLoading, isError) çeker.

```javascript
// ❌ YANLIŞ (Bileşen içinde doğrudan fetch)
function HomePage() {
  const [data, setData] = useState([]);
  useEffect(() => {
    fetch('/api/articles').then(res => res.json()).then(setData);
  }, []);
  return <div>...</div>;
}

//  DOĞRU (Custom Hook kullanımı)
import { useArticles } from '../../hooks/useArticles';

function HomePage() {
  const { data: articles, isLoading } = useArticles();
  if (isLoading) return <Skeleton />;
  return <ArticleList articles={articles} />;
}
```

---

## 🧪 4. Sıfır Console Hatası Standardı

Üretim ortamındaki deneyimi kusursuz kılmak için tarayıcı konsolunda **hiçbir uyarı, hata veya yakalanmamış istisna (uncaught exception)** kalmamalıdır.
* Bileşenlerin `prop-types` veya React 19 tip tanımlamaları eksiksiz olmalıdır.
* Listelerde render edilen her bileşene benzersiz ve kararlı bir `key` değeri verilmelidir (dizi indeksi kullanmaktan kaçının, `article.id` gibi benzersiz veriler tercih edin).
* Event listener'lar (örneğin window scroll) bileşen unmount edildiğinde temizlenmelidir (`useEffect` içindeki `return () => window.removeEventListener...` kalıbı).
