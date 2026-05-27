---
trigger: on_change
glob: frontend/src/**/*.{js,jsx}
description: HaberBot Durum Yönetimi, API ve B-Planı Entegrasyon Kuralları. TanStack React Query custom hooks, Promise.all paralel fetch ve Gemini API dayanıklılık fallbacks şablonları.
---

# 💾 HaberBot Durum Yönetimi, API & B-Planı Kuralları (rules/frontend/state-api.md)

Bu dosya, HaberBot uygulamasının ağ iletişimi, veri önbelleğe alma ve üçüncü taraf yapay zeka (Google Gemini) kotalarına bağlı hata senaryolarını yönetme standartlarını tanımlar. Custom hook yazarken, veri çekerken veya hata durumlarını ele alırken bu şablonlara sıkı sıkıya bağlı kalmalısınız.

*Bu dosya `vercel-react-best-practices` yeteneği standartlarına göre yapılandırılmıştır.*

---

## 🚀 1. TanStack React Query Custom Hook Standartları

HaberBot'ta veri çekme ve yönetme işlemleri tamamen `@tanstack/react-query` kütüphanesine emanettir. Bileşen içinde hiçbir zaman doğrudan state ataması ve fetch işlemi yapılmamalı; her fetch işlemi `services/api.js` içindeki API çağrısını sarmalayan **bir custom hook** olarak `hooks/` klasöründe tanımlanmalıdır.

### 🛠️ useQuery Custom Hook Şablonu (Makale Listeleme):
```javascript
import { useQuery } from '@tanstack/react-query';
import { fetchArticles } from '../services/api';

/**
 * Belirli bir konuya ait makaleleri çeken custom hook.
 * @param {string} topicId - Konu benzersiz kimliği
 */
export const useArticles = (topicId) => {
  return useQuery({
    queryKey: ['articles', topicId],
    queryFn: () => fetchArticles(topicId),
    staleTime: 1000 * 60 * 5,      // Veri 5 dakika boyunca taze (stale değil) kabul edilir.
    gcTime: 1000 * 60 * 10,       // Kullanılmayan cache 10 dakika sonra silinir (eski adıyla cacheTime).
    retry: 2,                     // Hata durumunda en fazla 2 kez tekrar dene.
    refetchOnWindowFocus: false,  // Tarayıcı sekmesi odağa geldiğinde otomatik tetiklemeyi kapat.
  });
};
```

### 🛠️ useMutation Custom Hook Şablonu (Yeni Kaynak Ekleme):
```javascript
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { addSource } from '../services/api';

export const useAddSource = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (newSourceData) => addSource(newSourceData),
    onSuccess: (data) => {
      // Başarılı eklemeden sonra 'sources' listesini otomatik geçersiz kıl (refetch tetikler).
      queryClient.invalidateQueries({ queryKey: ['sources'] });
    },
    onError: (error) => {
      console.error('Kaynak eklenirken bir hata oluştu:', error);
    }
  });
};
```

---

## 🌊 2. API Waterfall'ları Önleme (Paralel Fetching)

Bir bileşende veya sayfada birbirine bağımlı olmayan birden fazla veri çekme işlemi varsa, bunları sırayla (waterfall) beklemek yüklenme süresini katlar. Bunları **paralel** olarak başlatın.

```javascript
// ❌ YANLIŞ (Waterfall oluşur: topics bitmeden articles başlamaz)
const { data: topics } = useQuery({ queryKey: ['topics'], queryFn: fetchTopics });
const { data: articles } = useQuery({ queryKey: ['articles'], queryFn: fetchArticles });

//  DOĞRU (Custom hook'lar kendi içlerinde paralel çalışır, UI donmaz)
// React Query'de ardışık çağrılan hook'lar otomatik olarak paralel ağ istekleri atar.
```

Eğer bir API handler'ı içinde birden fazla bağımsız veri çekiliyorsa `Promise.all` kalıbını uygulayın:
```javascript
// services/api.js
export const fetchDashboardData = async () => {
  // Bağımsız API isteklerini paralel olarak koşturuyoruz.
  const [sources, stats, articles] = await Promise.all([
    fetchSources(),
    fetchStats(),
    fetchLatestArticles()
  ]);

  return { sources, stats, articles };
};
```

---

## 🛡️ 3. API Resilience & Gemini B-Planı (Fallback) Standartları

Backend'deki Google Gemini özetleme/çeviri servisi kota aşımı (`429 Too Many Requests`), internet kopması veya teknik bir arıza nedeniyle geçici olarak çalışmayabilir. Bu durumda web sitemiz **asla çökmeyecek, sonsuz yükleme döngüsünde kalmayacak ve boş/kırık veri göstermeyecektir**. 

Aşağıdaki yedek planları (fallbacks) uygulamak **mecburidir**:

### 🔤 Başlık Fallback Yapısı (Translated Title vs Original):
Gemini çevirisi gecikebilir veya başarısız olabilir. `title_tr` boşsa anında orijinal İngilizce başlığa dönün:
```javascript
// Bileşen İçinde Kullanım Şablonu
const displayTitle = article.title_tr || article.title || 'Başlıksız Makale';

return <h2 className="article-title">{displayTitle}</h2>;
```

### 📝 Özet ve Önizleme Fallback Yapısı (Summary vs Content):
Eğer Gemini henüz Türkçe özet (`summary_tr`) oluşturamadıysa, kullanıcıyı boş bırakmak yerine sırasıyla İngilizce özet (`summary`), ham İngilizce içerik (`original_content`) veya şık bir hazırlık mesajı gösterin:
```javascript
// Bileşen İçinde Kullanım Şablonu
const displaySummary = article.summary_tr 
  || article.summary 
  || (article.original_content ? `${article.original_content.substring(0, 150)}...` : null)
  || 'Özet ve Türkçe çeviri arka planda yapay zeka tarafından hazırlanıyor...';

return <p className="article-summary">{displaySummary}</p>;
```

### 📅 Güvenli Tarih Ayrıştırma (Safe Date Parsing):
Veritabanından veya backend modelinden gelen tarih alanları `created_at` veya `fetched_at` olarak farklı isimlendirmelere sahip olabilir. Bunları güvenli bir şekilde ayrıştırın:
```javascript
// Bileşen İçinde Kullanım Şablonu
export const formatArticleDate = (article) => {
  const targetDate = article.created_at || article.fetched_at;
  if (!targetDate) return 'Tarih Belirtilmedi';
  
  try {
    return new Date(targetDate).toLocaleDateString('tr-TR', {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  } catch (e) {
    return 'Geçersiz Tarih';
  }
};
```

---

## 🚫 4. State & API Anti-Patterns (Yapılmaması Gerekenler)

* **Sonsuz Yükleme Döngüsü Yaratmayın:** Arayüzde `isLoading` kontrolünü düzgün yapın. Ancak veri hataya düştüyse (`isError`) yükleniyor durumunda takılı kalmayın. Mutlaka bir `ErrorState` bileşeni render edin.
* **Türetilmiş Durumlar (Derived State) İçin Effect Kullanmayın:** Bir state'ten hesaplanabilen bir değişken için asla yeni bir `useState` ve `useEffect` açmayın. Render sırasında doğrudan hesaplayın:
  ```javascript
  // ❌ YANLIŞ
  const [activeArticles, setActiveArticles] = useState([]);
  useEffect(() => {
    setActiveArticles(articles.filter(a => a.active));
  }, [articles]);

  //  DOĞRU
  const activeArticles = articles.filter(a => a.active);
  ```
* **Query Key'leri String Olarak Bırakmayın:** Caching mekanizmasının doğru çalışması için query key'lerinizi yapılandırılmış diziler olarak tutun (örn: `['articles', topicId]`).
