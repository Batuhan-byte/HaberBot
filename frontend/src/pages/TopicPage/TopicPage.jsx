import React, { useState, useEffect, useMemo } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { api } from '../../services/api';
import { formatRelativeDate } from '../../utils/formatDate';
import Header from '../../components/Header/Header';
import './TopicPage.css';

// Fallback images representing modern tech portals
const FALLBACK_IMAGE_MAIN = "https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c?auto=format&fit=crop&w=800&q=80";
const FALLBACK_IMAGE_ALT = "https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?auto=format&fit=crop&w=600&q=80";

function stripHtmlTags(str) {
  if (!str) return '';
  return str.replace(/<\/?[^>]+(>|$)/g, "").trim();
}

/**
 * Skeletal loader for category grid.
 */
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

export default function TopicPage() {
  const { slug } = useParams();
  const [page, setPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState('');

  // Reset page when slug changes
  useEffect(() => {
    setPage(1);
  }, [slug]);

  // Fetch all topics (for navigation bar)
  const { data: topicsData } = useQuery({ 
    queryKey: ['topics'], 
    queryFn: api.getTopics 
  });
  
  const topics = topicsData?.topics || [];
  const currentTopic = topics.find(t => t.slug === slug);
  const topicName = currentTopic ? currentTopic.name : 'Teknoloji';
  
  // Fetch topic articles (paginated)
  const limit = 12;
  const { data: articlesData, isLoading, isError, refetch } = useQuery({
    queryKey: ['topicArticles', slug, page],
    queryFn: () => api.getTopicArticles(slug, page, limit),
    keepPreviousData: true
  });

  const articles = articlesData?.articles || [];
  const totalArticles = articlesData?.total || 0;
  const totalPages = Math.ceil(totalArticles / limit) || 1;

  // Extract Hero articles
  const gridMainArticle = articles[0];
  const gridSideArticles = useMemo(() => {
    return articles.slice(1, 5);
  }, [articles]);

  // Remaining articles for the bottom feed grid
  const feedArticles = useMemo(() => {
    return articles.slice(5);
  }, [articles]);

  // Dynamically compile a weekly AI category summary based on fetched articles
  const dynamicWeeklySummary = useMemo(() => {
    if (articles.length === 0) {
      return `Bu hafta ${topicName} kategorisinde henüz yeni bir içerik bulunmamaktadır. En güncel gelişmeler için takipte kalın.`;
    }
    const item1 = articles[0] ? `"${articles[0].title_tr || articles[0].title}"` : "";
    const item2 = articles[1] ? `"${articles[1].title_tr || articles[1].title}"` : "";
    const item3 = articles[2] ? `"${articles[2].title_tr || articles[2].title}"` : "";

    return `Bu hafta ${topicName} Haberleri ve İçerikleri kategorisinde teknoloji dünyası hareketli geçti. Özellikle öne çıkan ${item1} konusu okuyucular tarafından ilgiyle karşılandı. Ayrıca ${item2} ${item3 ? 've ' + item3 : ''} başlıkları en çok tartışılan ve değerlendirilen gelişmeler arasında yer aldı. HaberBot AI özetleme altyapısıyla hazırlanan bu derleme, kategorinin haftalık nabzını yansıtmaktadır.`;
  }, [articles, topicName]);

  const handlePageChange = (newPage) => {
    setPage(newPage);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

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
      <Header searchQuery={searchQuery} setSearchQuery={setSearchQuery} selectedTopic={topicName} />

      {/* ========== POPULER HASH-TAG BAR ========== */}
      <section className="populer-bar" aria-label="Popüler İçerikler">
        <span className="populer-title">Popüler İçerikler</span>
        <div className="populer-tags">
          {popularHashtags.map((tag, idx) => (
            <a key={idx} href={tag.link} className="tag-link">{tag.label}</a>
          ))}
        </div>
      </section>

      {/* ========== MAIN CONTENT CONTAINER ========== */}
      <main className="max-w-[1200px] mx-auto px-4 py-6">
        
        {/* Kategori Başlığı */}
        <div className="category-title-container mb-6">
          <h1 className="category-page-title">
            <span></span> {topicName} Haberleri ve İçerikleri
          </h1>
        </div>

        {/* HAFTANIN ÖZETİ (AI SUMMARY) PANELİ */}
        <section className="haftanin-ozeti-panel mb-8" aria-label="Haftanın Özeti">
          <div className="ozet-icon-wrap">
            <span className="material-symbols-outlined text-[#ff79c6] text-[24px]">insights</span>
          </div>
          <div className="ozet-content">
            <span className="ozet-badge">HAFTANIN ÖZETİ</span>
            <p className="ozet-text mt-2">{dynamicWeeklySummary}</p>
          </div>
        </section>

        {isError ? (
          <div className="text-center py-12 bg-gray-900 border border-gray-800 rounded-lg">
            <span className="material-symbols-outlined text-[48px] text-red-500 mb-2">error</span>
            <p className="text-gray-400">Haberler yüklenirken bir hata oluştu.</p>
            <button onClick={refetch} className="mt-4 px-6 py-2 bg-blue-600 hover:bg-blue-500 rounded text-sm">Tekrar Dene</button>
          </div>
        ) : isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {Array.from({ length: 6 }).map((_, i) => <ArticleSkeleton key={i} />)}
          </div>
        ) : articles.length === 0 ? (
          <div className="text-center py-12 bg-gray-900 border border-gray-800 rounded-lg">
            <span className="material-symbols-outlined text-[48px] text-gray-600 mb-4">search_off</span>
            <h3 className="text-lg font-bold text-white mb-2">Haber Bulunamadı</h3>
            <p className="text-gray-500">Bu kategoriye ait haber henüz bulunmamaktadır.</p>
          </div>
        ) : (
          <>
            {/* ========== CATEGORY HERO GRID (1-WIDE & 2x2 GRID) ========== */}
            <section className="category-hero-grid mb-8" aria-label="Öne Çıkanlar">
              {/* Sol Sütun: Büyük Manşet */}
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
                    <span className="cat-manset-badge">EDİTÖRÜN SEÇİMİ</span>
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
                        <span className="cat-sub-badge">{topicName.toUpperCase()}</span>
                        <h3 className="cat-sub-title">{displayTitle}</h3>
                      </div>
                    </Link>
                  );
                })}
              </div>
            </section>

            {/* ========== BOTTOM ARTICLES GRID ========== */}
            <section className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12" aria-label="Tüm Haberler">
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

          <div className="partner-logos-bar" aria-label="Marka İş Ortakları">
            <a href="#" className="partner-logo">mackolik</a>
            <a href="#" className="partner-logo">onedio</a>
            <a href="#" className="partner-logo">mynet</a>
            <a href="#" className="partner-logo">HaberBot</a>
            <a href="#" className="partner-logo">yemek.com</a>
            <a href="#" className="partner-logo">hisse.net</a>
          </div>

          <div className="text-center text-[0.7rem] text-gray-600 mt-6">
            © 2026 HaberBot AI & HaberBota Redesign. Tüm hakları saklıdır.
          </div>
        </div>
      </footer>
    </div>
  );
}
