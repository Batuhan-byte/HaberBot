import React, { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { api } from '../../services/api';
import { useDebounce } from '../../hooks/useDebounce';
import { formatRelativeDate } from '../../utils/formatDate';
import Header from '../../components/Header/Header';
import './HomePage.css';

// Fallback images representing modern tech portals
const FALLBACK_IMAGE_MAIN = "https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c?auto=format&fit=crop&w=800&q=80";
const FALLBACK_IMAGE_ALT = "https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?auto=format&fit=crop&w=600&q=80";

function stripHtmlTags(str) {
  if (!str) return '';
  return str.replace(/<\/?[^>]+(>|$)/g, "").trim();
}

/**
 * Skeletal loader for HaberBot row cards.
 */
function RowCardSkeleton() {
  return (
    <div className="row-card animate-pulse" aria-hidden="true">
      <div className="row-card-body">
        <div className="flex gap-2">
          <div className="h-4 w-12 bg-[#222222] rounded"></div>
          <div className="h-4 w-16 bg-[#222222] rounded"></div>
        </div>
        <div className="h-6 w-3/4 bg-[#222222] rounded mt-2"></div>
        <div className="space-y-1 mt-2">
          <div className="h-4 w-full bg-[#222222] rounded"></div>
          <div className="h-4 w-2/3 bg-[#222222] rounded"></div>
        </div>
        <div className="h-3 w-20 bg-[#222222] rounded mt-2"></div>
      </div>
      <div className="row-card-image-wrap bg-[#111111]"></div>
    </div>
  );
}

export default function HomePage() {
  const [selectedTopic, setSelectedTopic] = useState('Tümü');
  const [searchQuery, setSearchQuery] = useState('');
  const debouncedSearchQuery = useDebounce(searchQuery, 300);
  const [selectedSource, setSelectedSource] = useState('Tümü');

  // Fetch topics
  const { data: topics } = useQuery({ 
    queryKey: ['topics'], 
    queryFn: api.getTopics 
  });
  
  // Fetch recent articles
  const { data: recentArticles, isLoading: recentLoading } = useQuery({ 
    queryKey: ['articles'], 
    queryFn: api.getArticles 
  });

  // Fetch search results if search is active
  const isSearchActive = debouncedSearchQuery.length > 0 || selectedSource !== 'Tümü';
  const { data: searchResults, isLoading: searchLoading } = useQuery({
    queryKey: ['search', debouncedSearchQuery, selectedSource],
    queryFn: () => api.searchArticles(debouncedSearchQuery, selectedSource),
    enabled: isSearchActive
  });

  const topicsList = topics?.topics && Array.isArray(topics.topics) ? topics.topics : [];
  
  const rawArticlesList = isSearchActive 
    ? (searchResults?.articles || []) 
    : (recentArticles?.articles || []);

  const articlesLoading = isSearchActive ? searchLoading : recentLoading;

  // Filter articles based on selected topic
  const filteredArticles = useMemo(() => {
    let list = Array.isArray(rawArticlesList) ? rawArticlesList : [];
    if (selectedTopic !== 'Tümü') {
      const topic = topicsList.find(t => t.name === selectedTopic);
      if (topic) {
        list = list.filter(article => article.topic_id === topic.id);
      }
    }
    return list;
  }, [rawArticlesList, topicsList, selectedTopic]);

  // Sort articles by score (descending) to find the most read/popular trend items
  const popularArticles = useMemo(() => {
    const list = Array.isArray(rawArticlesList) ? [...rawArticlesList] : [];
    return list.sort((a, b) => (b.score || 0) - (a.score || 0)).slice(0, 5);
  }, [rawArticlesList]);

  // Extract Hero Articles
  const topHeroArticle = filteredArticles[0];
  const secondHeroArticle = filteredArticles[1];
  
  // Remaining articles for the list akış
  const feedArticles = useMemo(() => {
    return filteredArticles.slice(2);
  }, [filteredArticles]);

  // Otomobil specific lane articles (for category showcase)
  const otomobilArticles = useMemo(() => {
    const list = Array.isArray(rawArticlesList) ? rawArticlesList : [];
    const autoTopic = topicsList.find(t => t.slug === 'otomobil' || t.name === 'Otomobil');
    if (autoTopic) {
      return list.filter(article => article.topic_id === autoTopic.id).slice(0, 4);
    }
    return list.slice(0, 4); // fallback if otomobil doesn't exist
  }, [rawArticlesList, topicsList]);

  // Static indirim kuponları lists (Sıcak Fırsatlar)
  const firsatlar = [
    { title: "Amazon'da Günün En İyi Teknoloji Fırsatları", code: "Sıcak Fırsat", link: "#" },
    { title: "Tüm Adidas Ayakkabılarda HaberBot30 ile %30 İndirim!", code: "İndirim Kodu", link: "#" },
    { title: "Steam İlkbahar İndirimlerinde Kaçırılmayacak 5 Oyun", code: "Haber", link: "#" },
    { title: "Protein Ocean Siparişlerinde HaberBot Koduyla %10 İndirim", code: "Sponsorlu", link: "#" }
  ];

  // Static Video items
  const videos = [
    { title: "Huawei Watch Fit 5 Pro İncelemesi: Fiyatı En İddialı Akıllı Saat!", image: "https://images.unsplash.com/photo-1544111892-4a9606cb0b82?auto=format&fit=crop&w=300&q=80", link: "#" },
    { title: "Ustalık Eseri Fiyat Performans Telefonu | Honor 600 İnceleme", image: "https://images.unsplash.com/photo-1598327105666-5b89351aff97?auto=format&fit=crop&w=300&q=80", link: "#" },
    { title: "Bit Pazarı Elektronik Aletlerini İndirim Kuponuna Dönüştürdük", image: "https://images.unsplash.com/photo-1531403009284-440f080d1e12?auto=format&fit=crop&w=300&q=80", link: "#" },
    { title: "Mark Zuckerberg Kendini Klonlayacak | Teknoloji Gündemi", image: "https://images.unsplash.com/photo-1563986768609-322da13575f3?auto=format&fit=crop&w=300&q=80", link: "#" },
    { title: "Robotik Koluyla Dip Köşe Temizleyen Homend Alex 50 İnceleme", image: "https://images.unsplash.com/photo-1558317374-067fb5f30001?auto=format&fit=crop&w=300&q=80", link: "#" }
  ];

  const popularHashtags = [
    { label: "#Google/O2026", link: "#" },
    { label: "#Apple", link: "#" },
    { label: "#Samsung", link: "#" },
    { label: "#Xiaomi", link: "#" },
    { label: "#ElektrikliOtomobil", link: "#" },
    { label: "#ChatGPT", link: "#" },
    { label: "#OyunIndirimleri", link: "#" },
    { label: "#UcretsizOyun", link: "#" }
  ];

  return (
    <div className="min-h-screen bg-black text-white selection:bg-blue-600 selection:text-white">
      {/* ========== WEBTEKNO NAVIGATION BAR ========== */}
      <Header searchQuery={searchQuery} setSearchQuery={setSearchQuery} selectedTopic="Tümü" />

      {/* ========== POPULER HASH-TAG BAR ========== */}
      <section className="populer-bar" aria-label="Popüler İçerikler">
        <span className="populer-title">Popüler İçerikler</span>
        <div className="populer-tags">
          {popularHashtags.map((tag, idx) => (
            <a key={idx} href={tag.link} className="tag-link">{tag.label}</a>
          ))}
        </div>
      </section>

      {/* ========== WEBTEKNO UPPER 3-COLUMN HERO GRID (REKAMLAR ELENDİ) ========== */}
      <section className="HaberBot-hero-grid" aria-label="Manşet Haberler">
        {/* Sol Sütun: Sıcak Fırsatlar */}
        <div className="firsat-box">
          <h2 className="firsat-header">
            <span className="material-symbols-outlined text-[18px]">local_mall</span>
            Sıcak Fırsatlar
          </h2>
          {firsatlar.map((f, i) => (
            <a href={f.link} key={i} className="firsat-item">
              <span className="uppercase tracking-wider font-bold">{f.code}</span>
              <p className="line-clamp-2">{f.title}</p>
            </a>
          ))}
        </div>

        {/* Orta Sütun: Büyük Manşet */}
        {topHeroArticle ? (
          <Link to={`/haber/${topHeroArticle.id}`} className="manset-card">
            <div className="manset-image-wrap">
              <img 
                src={topHeroArticle.image_url || FALLBACK_IMAGE_MAIN} 
                alt={topHeroArticle.title} 
                className="manset-image" 
              />
              <div className="manset-overlay"></div>
            </div>
            <div className="manset-body">
              <span className="manset-badge">MANŞET</span>
              <h2 className="manset-title">{topHeroArticle.title_tr || topHeroArticle.title}</h2>
              <span className="manset-subtitle">
                {stripHtmlTags(topHeroArticle.turkish_summary || topHeroArticle.summary || topHeroArticle.original_content).slice(0, 60)}...
              </span>
            </div>
          </Link>
        ) : (
          <div className="manset-card flex items-center justify-center bg-gray-900 border border-gray-800">
            <div className="text-center p-8">
              <span className="material-symbols-outlined text-[48px] text-gray-600 mb-2">article</span>
              <p className="text-gray-500">Haber yükleniyor...</p>
            </div>
          </div>
        )}

        {/* Sağ Sütun: Düz Vurgu Kartı */}
        {secondHeroArticle ? (
          <Link to={`/haber/${secondHeroArticle.id}`} className="vurgu-card">
            <div>
              <span className="vurgu-badge">GÜNDEM</span>
              <h2 className="vurgu-title mt-4">{secondHeroArticle.title_tr || secondHeroArticle.title}</h2>
            </div>
            <p className="vurgu-desc">
              {stripHtmlTags(secondHeroArticle.turkish_summary || secondHeroArticle.summary || secondHeroArticle.original_content)}
            </p>
          </Link>
        ) : (
          <div className="vurgu-card flex items-center justify-center bg-gray-900 border border-gray-800">
            <p className="text-gray-500 text-sm">Haber bekleniyor...</p>
          </div>
        )}
      </section>

      {/* ========== POPULER VIDEOLAR SLIDER ========== */}
      <section className="video-section" aria-label="Popüler Videolar">
        <h2 className="section-header">
          <span>▶</span> Popüler Videolar
        </h2>
        <div className="video-slider">
          {videos.map((vid, i) => (
            <a href={vid.link} key={i} className="video-card">
              <div className="video-thumb-wrap">
                <img src={vid.image} alt={vid.title} className="video-thumb" />
                <div className="play-btn-overlay">
                  <div className="play-icon">▶</div>
                </div>
              </div>
              <p className="video-title">{vid.title}</p>
            </a>
          ))}
        </div>
      </section>

      {/* ========== TWO-COLUMN MAIN CONTENT & SIDEBAR ========== */}
      <main className="HaberBot-main-layout">
        {/* Sol Sütun: Dikey Yatay Haber Akışı */}
        <section className="news-feed-list" aria-label="Teknoloji Haberleri">
          {articlesLoading ? (
            Array(4).fill(0).map((_, i) => <RowCardSkeleton key={i} />)
          ) : feedArticles && feedArticles.length > 0 ? (
            feedArticles.map((article) => {
              const displayTitle = article.title_tr || article.title;
              const cleanSummary = stripHtmlTags(article.turkish_summary || article.summary || article.original_content || 'Detaylar yükleniyor...');
              
              return (
                <Link to={`/haber/${article.id}`} key={article.id} className="row-card">
                  <div className="row-card-body">
                    <div className="row-card-meta">
                      <span className="row-card-badge">{article.source === 'hackernews' ? 'HN' : 'RSS'}</span>
                      <span>{formatRelativeDate(article.created_at || article.fetched_at)}</span>
                    </div>
                    <h3 className="row-card-title">{displayTitle}</h3>
                    <p className="row-card-preview">{cleanSummary}</p>
                    <span className="row-card-author">Deniz Şen — 3 saat önce</span>
                  </div>
                  <div className="row-card-image-wrap">
                    <img 
                      src={article.image_url || FALLBACK_IMAGE_ALT} 
                      alt={displayTitle} 
                      className="row-card-image"
                    />
                  </div>
                </Link>
              );
            })
          ) : (
            <div className="py-12 flex flex-col items-center justify-center text-center bg-gray-900 border border-gray-800 rounded-lg p-8">
              <span className="material-symbols-outlined text-[48px] text-gray-700 mb-4">search_off</span>
              <p className="text-gray-500 font-body-lg text-body-lg">
                {isSearchActive ? 'Arama kriterlerinize uygun haber bulunamadı.' : 'Bu kategoriye ait haber henüz bulunmamaktadır.'}
              </p>
            </div>
          )}
        </section>

        {/* Sağ Sütun: Sidebar Widget'ları */}
        <aside className="HaberBot-sidebar" aria-label="Yan Menü">
          <div className="trend-box">
            <h2 className="trend-title">
              <span>★</span> En Çok Okunanlar
            </h2>
            
            {/* Top Trend Item */}
            {popularArticles[0] && (
              <Link to={`/haber/${popularArticles[0].id}`} className="trend-top-card">
                <img 
                  src={popularArticles[0].image_url || FALLBACK_IMAGE_ALT} 
                  alt={popularArticles[0].title}
                  className="absolute inset-0 w-full h-full object-cover" 
                />
                <div className="trend-top-overlay"></div>
                <div className="trend-top-body">
                  <span className="trend-top-badge">TREND 1</span>
                  <h3 className="trend-top-title">{popularArticles[0].title_tr || popularArticles[0].title}</h3>
                </div>
              </Link>
            )}

            {/* Trend list items 2-5 */}
            <div className="trend-list">
              {popularArticles.slice(1).map((article, idx) => (
                <Link to={`/haber/${article.id}`} key={article.id} className="trend-item">
                  <span className="trend-number">{idx + 2}</span>
                  <div className="trend-item-body">
                    <h4 className="trend-item-title">{article.title_tr || article.title}</h4>
                    <span className="trend-item-source">{article.source || 'KAYNAK'}</span>
                  </div>
                </Link>
              ))}
            </div>
          </div>
        </aside>
      </main>

      {/* ========== KATEGORİSEL YATAY AKIŞ (OTOMOBİL LANESI) ========== */}
      <section className="category-lane" aria-label="Kategorisel Akış">
        <div className="lane-header">
          <h2 className="lane-title">
            <span>🚗</span> Otomobil Haberleri
          </h2>
        </div>
        <div className="lane-slider">
          {otomobilArticles.map((article) => {
            const displayTitle = article.title_tr || article.title;
            return (
              <Link to={`/haber/${article.id}`} key={article.id} className="lane-card">
                <div className="lane-card-image-wrap">
                  <img 
                    src={article.image_url || FALLBACK_IMAGE_ALT} 
                    alt={displayTitle} 
                    className="lane-card-image"
                  />
                </div>
                <div className="lane-card-body">
                  <span className="lane-card-badge">OTOMOBİL</span>
                  <h3 className="lane-card-title">{displayTitle}</h3>
                </div>
              </Link>
            );
          })}
        </div>
      </section>

      {/* ========== WEBTEKNO SITEMAP FOOTER ========== */}
      <footer className="sitemap-footer" aria-label="Sayfa Alt Bilgisi">
        <div className="sitemap-container">
          <div className="sitemap-grid">
            <div className="sitemap-brand-col">
              <span className="sitemap-brand-logo">
                Haber<span>Bot</span>
              </span>
              <p className="sitemap-brand-desc">
                Türkiye'nin en popüler teknoloji haber ve inceleme platformu. HaberBot AI altyapısıyla çalışır.
              </p>
              <div className="sitemap-socials">
                <a href="#" className="sitemap-social-icon material-symbols-outlined" aria-label="Facebook">public</a>
                <a href="#" className="sitemap-social-icon material-symbols-outlined" aria-label="Twitter">terminal</a>
                <a href="#" className="sitemap-social-icon material-symbols-outlined" aria-label="Instagram">photo_camera</a>
                <a href="#" className="sitemap-social-icon material-symbols-outlined" aria-label="YouTube">smart_display</a>
              </div>
            </div>

            <div className="sitemap-col">
              <h3 className="sitemap-heading">Kurumsal</h3>
              <div className="sitemap-links">
                <a href="#" className="sitemap-link">Hakkımızda</a>
                <a href="#" className="sitemap-link">Künye</a>
                <a href="#" className="sitemap-link">İletişim</a>
                <a href="#" className="sitemap-link">Gizlilik Sözleşmesi</a>
              </div>
            </div>

            <div className="sitemap-col">
              <h3 className="sitemap-heading">Kategoriler</h3>
              <div className="sitemap-links">
                <a href="#" className="sitemap-link">Yapay Zeka</a>
                <a href="#" className="sitemap-link">Mobil</a>
                <a href="#" className="sitemap-link">Oyun</a>
                <a href="#" className="sitemap-link">Otomobil</a>
              </div>
            </div>

            <div className="sitemap-col">
              <h3 className="sitemap-heading">En Çok Paylaşılanlar</h3>
              <div className="sitemap-links">
                <a href="#" className="sitemap-link">Haftanın Haberleri</a>
                <a href="#" className="sitemap-link">Hacker News Trendleri</a>
                <a href="#" className="sitemap-link">Teknoloji Gündemi</a>
              </div>
            </div>

            <div className="sitemap-col">
              <h3 className="sitemap-heading">En Çok İzlenenler</h3>
              <div className="sitemap-links">
                <a href="#" className="sitemap-link">Ürün İncelemeleri</a>
                <a href="#" className="sitemap-link">Akıllı Telefon Testleri</a>
                <a href="#" className="sitemap-link">Yapay Zeka Karşılaştırmaları</a>
              </div>
            </div>
          </div>

          {/* Partner Brands */}
          <div className="partner-logos-bar" aria-label="Marka İş Ortakları">
            <a href="#" className="partner-logo">mackolik</a>
            <a href="#" className="partner-logo">onedio</a>
            <a href="#" className="partner-logo">mynet</a>
            <a href="#" className="partner-logo">HaberBot</a>
            <a href="#" className="partner-logo">yemek.com</a>
            <a href="#" className="partner-logo">hisse.net</a>
          </div>

          <div className="text-center text-[0.7rem] text-gray-600 mt-6">
            © 2026 HaberBot AI & HaberBot Redesign. Tüm hakları saklıdır.
          </div>
        </div>
      </footer>
    </div>
  );
}
