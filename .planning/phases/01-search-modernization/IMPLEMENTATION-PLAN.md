# Implementasyon Planı: Gelişmiş Arama Sistemi

**Tarih:** 2026-01-27
**Versiyon:** 1.0
**Durum:** Planlanıyor

---

## Özet

Mevcut arama sistemi iyileştirilerek:
1. Yazarken dropdownda 3 öneri haber gösterilecek
2. Entera basınca ayrı bir sonuç sayfasına yönlendirilecek
3. Sonuç sayfası TopicPage düzeninde olacak

---

## Mevcut Durum Analizi

### Mimari Akış
[Kullanıcı] -> Header.jsx (localSearch) -> HomePage (setSearchQuery) -> useDebounce(300ms) -> API -> Backend -> SQL (ILIKE) -> Sonuçlar

### Değişiklik Gereken Noktalar
| Katman | Dosya | Değişiklik |
|--------|-------|------------|
| Backend | postgres_article.go | Relevance scoring ekle |
| Backend | article_handler.go | Yeni endpoint (limit=3) |
| Frontend | SearchDropdown.jsx | Yeni bileşen |
| Frontend | Header.jsx | Dropdown entegrasyonu |
| Frontend | Header.css | Dropdown stilleri |
| Frontend | SearchPage.jsx | Yeni sayfa |
| Frontend | App.jsx | Route ekleme |
| Frontend | api.js | Limit=3 API çağrısı |

---

## Faz 01: Backend - Relevance Sıralama

**Hedef:** Arama sonuçlarını en alakalı olana göre sırala

### 1.1 SQL Query Güncelleme

**Dosya:** backend/internal/adapter/repository/postgres_article.go

**Mevcut (satır ~151, ~164):**
ORDER BY fetched_at DESC

**Yeni:**
ORDER BY 
    CASE 
        WHEN title ILIKE ''%'' || $1 || ''%'' THEN 3
        WHEN turkish_title ILIKE ''%'' || $1 || ''%'' THEN 2
        WHEN turkish_summary ILIKE ''%'' || $1 || ''%'' THEN 1
        ELSE 0
    END DESC,
    fetched_at DESC

---

## Faz 02: Frontend - SearchDropdown Bileşeni

**Hedef:** 2+ karakter yazılınca 3 öneri göster

### 2.1 Dosya Yapısı
frontend/src/components/SearchDropdown/
  SearchDropdown.jsx
  SearchDropdown.css

### 2.2 Bileşen Şablonu

const SearchDropdown = ({ query, onSelect }) => {
    const { data, isLoading } = useQuery({
        queryKey: [''search-suggestions'', query],
        queryFn: () => api.searchArticles(query, '''', 1, 3),
        enabled: query.length >= 2
    });

    if (!data?.articles?.length) return null;

    return (
        <div className="search-dropdown">
            {data.articles.map(article => (
                <div key={article.id} className="dropdown-item" onClick={() => onSelect(article)}>
                    <img src={article.image_url} alt="" className="dropdown-thumb" />
                    <span className="dropdown-title">{article.turkish_title || article.title}</span>
                </div>
            ))}
        </div>
    );
};

### 2.3 CSS Stilleri

.search-dropdown {
    position: absolute;
    top: 100%;
    right: 0;
    margin-top: 8px;
    width: 320px;
    background: #111111;
    border: 1px solid rgba(255,255,255,0.1);
    border-radius: 12px;
    z-index: 1000;
    box-shadow: 0 10px 40px rgba(0,0,0,0.5);
}

.dropdown-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    cursor: pointer;
}

.dropdown-item:hover {
    background: rgba(59, 130, 246, 0.1);
}

.dropdown-thumb {
    width: 40px;
    height: 40px;
    object-fit: cover;
    border-radius: 6px;
}

---

## Faz 03: Header Entegrasyonu

**Hedef:** Arama kutusuna dropdown ekle ve Enter davranışını değiştir

### 3.1 State Eklenecek
const [showDropdown, setShowDropdown] = useState(false);
const debouncedSearch = useDebounce(localSearch, 200);

### 3.2 Enter Handler
const handleSearchSubmit = (e) => {
    if (e.key === ''Enter'' && localSearch.trim()) {
        navigate(`/arama?q=${encodeURIComponent(localSearch)}`);
    }
};

### 3.3 Dropdown Entegrasyonu
{showDropdown && localSearch.length >= 2 && (
    <SearchDropdown 
        query={debouncedSearch}
        onSelect={(article) => {
            navigate(`/haber/${article.id}`);
            setLocalSearch('');
            setShowDropdown(false);
        }}
    />
)}

---

## Faz 04: Arama Sonuç Sayfası

**Hedef:** /arama?q=... route oluştur

### 4.1 App.jsx
<Route path="/arama" element={<SearchPage />} />

### 4.2 SearchPage.jsx
const SearchPage = () => {
    const [searchParams] = useSearchParams();
    const query = searchParams.get(''q'') || '''';
    
    const { data, isLoading } = useQuery({
        queryKey: [''search'', query],
        queryFn: () => api.searchArticles(query),
        enabled: !!query
    });

    return (
        <div className="search-page">
            <Header />
            <main>
                <h1>Sonuç: "{query}"</h1>
                {/* Article grid */}
            </main>
        </div>
    );
};

---

## Görev Listesi

| # | Görev | Dosyalar | Durum |
|---|-------|----------|-------|
| 1 | SQL Relevance Scoring | postgres_article.go | Beklemede |
| 2 | SearchDropdown Bileşeni | SearchDropdown.jsx, .css | Beklemede |
| 3 | Header Entegrasyonu | Header.jsx, .css | Beklemede |
| 4 | SearchPage Oluşturma | SearchPage.jsx, .css, App.jsx | Beklemede |
| 5 | Test ve Doğrulama | - | Beklemede |

---

## Success Criteria

- [ ] 2+ karakter yazılınca dropdown açılıyor
- [ ] Dropdownda 3 öneri görünüyor
- [ ] Dropdownda tıklama haber sayfasına gidiyor
- [ ] Entera basınca /arama?q=... sayfasına gidiyor
- [ ] Arama sonuç sayfası TopicPage düzeninde
- [ ] Build başarılı

---

## Bağımlılıklar

Görev 1 (Backend) -> paralel -> Görev 2 (Dropdown)
                              |
                        Görev 3 (Header)
                              |
                        Görev 4 (SearchPage)
                              |
                        Görev 5 (Test)

---

## Notlar

- Mevcut useDebounce hooku zaten var
- TopicPage.jsx referans alınarak SearchPage oluşturulacak