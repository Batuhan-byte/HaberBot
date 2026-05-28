import React, { useState, useEffect, useMemo } from 'react';
import { useSearchParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { api } from '../../services/api';
import { formatRelativeDate } from '../../utils/formatDate';
import Header from '../../components/Header/Header';
import './SearchPage.css';

const FALLBACK_IMAGE_MAIN = "https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c?auto=format&fit=crop&w=800&q=80";
const FALLBACK_IMAGE_ALT = "https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?auto=format&fit=crop&w=600&q=80";

function stripHtmlTags(str) {
  if (!str) return '';
  return str.replace(/<\/?[^>]+(>|$)/g, "").trim();
}

function ArticleSkeleton() {
  return (
    <div className="bg-[#0A0A0A] border border-[#222222] p-6 rounded-lg flex flex-col gap-4 animate-pulse" aria-hidden="true">
      <div className="h-4 w-20 bg-[#222222] rounded"></div>
      <div className="h-6 w-full bg-[#222222] rounded mt-2"></div>
      <div className="space-y-2 mt-2">
        <div className="h-4 w-full bg-[#222222] rounded"></div>
        <div className="h-4 w-2/3 bg-[#222222] rounded"></div>
      </div>
    </div>
  );
}

export default function SearchPage() {
  const [searchParams] = useSearchParams();
  const query = searchParams.get('q') || '';
  const [page, setPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState('');

  // Reset page when search query parameter changes
  useEffect(() => {
    setPage(1);
  }, [query]);

  const limit = 12;
  
  // Fetch paginated search results
  const { data: searchResults, isLoading, isError, refetch } = useQuery({
    queryKey: ['search-results', query, page],
    queryFn: () => api.searchArticles(query, '', page, limit),
    enabled: !!query,
    keepPreviousData: true
  });

  const articles = searchResults?.articles || [];
  const totalArticles = searchResults?.total || 0;
  const totalPages = Math.ceil(totalArticles / limit) || 1;

  // Extract Hero articles for TopicPage layout visual consistency
  const gridMainArticle = articles[0];
  const gridSideArticles = useMemo(() => {
    return articles.slice(1, 5);
  }, [articles]);

  const feedArticles = useMemo(() => {
    return articles.slice(5);
  }, [articles]);

  const handlePageChange = (newPage) => {
    setPage(newPage);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const popularHashtags = [
    { label: "#Google/O2026", link: "/arama?q=Google" },
    { label: "#Apple", link: "/arama?q=Apple" },
    { label: "#Samsung", link: "/arama?q=Samsung" },
    { label: "#Xiaomi", link: "/arama?q=Xiaomi" },
    { label: "#ElektrikliOtomobil", link: "/arama?q=Elektrikli" },
    { label: "#ChatGPT", link: "/arama?q=ChatGPT" },
    { label: "#OyunIndirimleri", link: "/arama?q=Oyun" },
    { label: "#UcretsizOyun", link: "/arama?q=Bedava" }
  ];

  return (
    <div className="min-h-screen bg-black text-white selection:bg-blue-600 selection:text-white">
      {/* ========== HABERBOT NAVIGATION BAR ========== */}
      <Header searchQuery={searchQuery} setSearchQuery={setSearchQuery} />

      {/* ========== POPULER HASH-TAG BAR ========== */}
      <section className="populer-bar" aria-label="Popüler İçerikler">
        <span className="populer-title">Popüler İçerikler</span>
        <div className="populer-tags">
          {popularHashtags.map((tag, idx) => (
            <Link key={idx} to={tag.link} className="tag-link">{tag.label}</Link>
          ))}
        </div>
      </section>

      {/* ========== MAIN CONTENT CONTAINER ========== */}
      <main className="max-w-[1200px] mx-auto px-4 py-6">
        
        {/* Arama Sonuç Başlığı */}
        <div className="search-title-container mb-6">
          <h1 className="search-page-title">
            <span></span> Arama Sonuçları: "{query}"
            <span className="search-count-badge">{totalArticles} Haber</span>
          </h1>
        </div>

        {isError ? (
          <div className="text-center py-12 bg-gray-900 border border-gray-800 rounded-lg">
            <span className="material-symbols-outlined text-[48px] text-red-500 mb-2">error</span>
            <p className="text-gray-400">Arama sonuçları yüklenirken bir hata oluştu.</p>
            <button onClick={refetch} className="mt-4 px-6 py-2 bg-blue-600 hover:bg-blue-500 rounded text-sm font-bold">Tekrar Dene</button>
          </div>
        ) : isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {Array.from({ length: 6 }).map((_, i) => <ArticleSkeleton key={i} />)}
          </div>
        ) : articles.length === 0 ? (
          <div className="text-center py-16 bg-gray-900 border border-gray-800 rounded-lg hb-search-no-results">
            <span className="material-symbols-outlined text-[64px] text-gray-600 mb-4">search_off</span>
            <h3 className="text-xl font-bold text-white mb-2">Aramanızla Eşleşen Haber Bulunamadı</h3>
            <p className="text-gray-400 max-w-md mx-auto mb-6 text-sm">
              Lütfen farklı anahtar kelimeler deneyin veya yukarıdaki popüler etiketleri kullanarak arama yapın.
            </p>
            <div className="flex flex-wrap justify-center gap-2 max-w-lg mx-auto">
              {popularHashtags.map((tag, idx) => (
                <Link key={idx} to={tag.link} className="px-3 py-1.5 bg-black hover:bg-zinc-900 border border-zinc-800 rounded-full text-xs font-bold text-gray-300 transition-colors">
                  {tag.label}
                </Link>
              ))}
            </div>
          </div>
        ) : (
          <>
            {/* ========== TOPICPAGE STYLE HERO GRID (1-WIDE & 2x2 GRID) ========== */}
            <section className="category-hero-grid mb-8" aria-label="Öne Çıkan Arama Sonuçları">
              {/* Sol Sütun: En Alakalı Büyük Manşet */}
              {gridMainArticle && (
                <Link to={`/haber/${gridMainArticle.id}`} className="cat-manset-card">
                  <div className="cat-manset-image-wrap">
                    <img 
                      src={gridMainArticle.image_url || FALLBACK_IMAGE_MAIN} 
                      alt={gridMainArticle.title} 
                      className="cat-manset-image"
                    />
                    <div className="cat-manset-overlay"></div>
                  </div>
                  <div className="cat-manset-body">
                    <span className="cat-manset-badge">EN ALAKALI SONUÇ</span>
                    <h2 className="cat-manset-title">{gridMainArticle.title_tr || gridMainArticle.title}</h2>
                  </div>
                </Link>
              )}

              {/* Sağ Sütun: 2x2 Grid */}
              <div className="cat-grid-side">
                {gridSideArticles.map((article) => {
                  const displayTitle = article.title_tr || article.title;
                  return (
                    <Link to={`/haber/${article.id}`} key={article.id} className="cat-sub-card">
                      <div className="cat-sub-image-wrap">
                        <img 
                          src={article.image_url || FALLBACK_IMAGE_ALT} 
                          alt={displayTitle} 
                          className="cat-sub-image"
                        />
                      </div>
                      <div className="cat-sub-body">
                        <span className="cat-sub-badge">BENZER SONUÇ</span>
                        <h3 className="cat-sub-title">{displayTitle}</h3>
                      </div>
                    </Link>
                  );
                })}
              </div>
            </section>

            {/* ========== BOTTOM FEED GRID ========== */}
            <section className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12" aria-label="Diğer Arama Sonuçları">
              {feedArticles.map((article) => {
                const displayTitle = article.title_tr || article.title;
                const cleanSummary = stripHtmlTags(article.turkish_summary || article.summary || article.original_content).slice(0, 100);
                return (
                  <Link to={`/haber/${article.id}`} key={article.id} className="bg-[#0d0d0d] border border-[#222222] hover:border-[#3b82f6] rounded-lg overflow-hidden flex flex-col transition-all duration-200 hover:-translate-y-1">
                    <div className="h-44 bg-gray-950 overflow-hidden">
                      <img 
                        src={article.image_url || FALLBACK_IMAGE_ALT} 
                        alt={displayTitle} 
                        className="w-full h-full object-cover transition-transform duration-300 hover:scale-103"
                      />
                    </div>
                    <div className="p-4 flex flex-col flex-1 justify-between gap-4">
                      <div>
                        <div className="flex items-center gap-2 text-[0.7rem] text-gray-500 mb-2">
                          <span className="bg-[#222222] text-gray-300 px-2 py-0.5 rounded uppercase font-bold">{article.source === 'hackernews' ? 'HN' : 'RSS'}</span>
                          <span>{formatRelativeDate(article.created_at || article.fetched_at)}</span>
                        </div>
                        <h3 className="font-bold text-white leading-snug line-clamp-2 hover:text-[#3b82f6] transition-colors">{displayTitle}</h3>
                        <p className="text-gray-400 text-xs line-clamp-3 mt-2">{cleanSummary}...</p>
                      </div>
                      <span className="text-[0.7rem] text-gray-600 font-semibold mt-auto">Deniz Şen — {formatRelativeDate(article.created_at || article.fetched_at)}</span>
                    </div>
                  </Link>
                );
              })}
            </section>

            {/* ========== PAGINATION ========== */}
            {totalPages > 1 && (
              <div className="flex justify-center mt-12 gap-2">
                <button 
                  disabled={page === 1}
                  onClick={() => handlePageChange(page - 1)}
                  className="px-4 py-2 border border-[#222222] hover:border-[#3b82f6] disabled:opacity-40 disabled:hover:border-[#222222] rounded transition-all text-xs font-bold"
                >
                  Önceki
                </button>
                <span className="px-4 py-2 text-xs font-bold bg-[#111111] border border-[#222222] rounded text-[#3b82f6]">
                  {page} / {totalPages}
                </span>
                <button 
                  disabled={page === totalPages}
                  onClick={() => handlePageChange(page + 1)}
                  className="px-4 py-2 border border-[#222222] hover:border-[#3b82f6] disabled:opacity-40 disabled:hover:border-[#222222] rounded transition-all text-xs font-bold"
                >
                  Sonraki
                </button>
              </div>
            )}
          </>
        )}
      </main>

      {/* ========== HABERBOT SITEMAP FOOTER ========== */}
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