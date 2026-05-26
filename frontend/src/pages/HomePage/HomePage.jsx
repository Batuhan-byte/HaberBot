import React, { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { api } from '../../services/api';
import ArticleCard from '../../components/ArticleCard/ArticleCard';
import { useDebounce } from '../../hooks/useDebounce';

/**
 * Skeleton matching the Terminal Monolith card structure.
 */
function ArticleSkeleton() {
  return (
    <div className="bg-[#0A0A0A] border border-[#222222] p-6 rounded-lg flex flex-col gap-4 animate-pulse" aria-hidden="true">
      <div className="flex justify-between items-start">
        <div className="h-5 w-24 bg-[#111111] rounded"></div>
        <div className="h-4 w-20 bg-[#111111] rounded"></div>
      </div>
      <div className="h-7 w-full bg-[#111111] rounded mt-1"></div>
      <div className="space-y-2 mt-1">
        <div className="h-4 w-full bg-[#111111] rounded"></div>
        <div className="h-4 w-full bg-[#111111] rounded"></div>
        <div className="h-4 w-2/3 bg-[#111111] rounded"></div>
      </div>
      <div className="mt-auto pt-4 border-t border-[#111111] flex gap-2">
        <div className="h-3 w-10 bg-[#111111] rounded"></div>
        <div className="h-3 w-10 bg-[#111111] rounded"></div>
      </div>
    </div>
  );
}

export default function HomePage() {
  const [selectedTopic, setSelectedTopic] = useState('Tümü');
  
  // Arama state'i
  const [searchQuery, setSearchQuery] = useState('');
  // Performans için 300ms debounce (Kullanıcı yazmayı bitirene kadar bekler, gereksiz API çağrılarını önler)
  const debouncedSearchQuery = useDebounce(searchQuery, 300);
  
  // Kaynak (Source) filtresi state'i
  const [selectedSource, setSelectedSource] = useState('Tümü');

  // Konuları getir
  const { data: topics } = useQuery({ 
    queryKey: ['topics'], 
    queryFn: api.getTopics 
  });
  
  // Normal akış için en son makaleler
  const { data: recentArticles, isLoading: recentLoading } = useQuery({ 
    queryKey: ['articles'], 
    queryFn: api.getArticles 
  });

  // Arama veya kaynak filtresi aktifse Search API'sini çağır
  const isSearchActive = debouncedSearchQuery.length > 0 || selectedSource !== 'Tümü';
  const { data: searchResults, isLoading: searchLoading } = useQuery({
    queryKey: ['search', debouncedSearchQuery, selectedSource],
    queryFn: () => api.searchArticles(debouncedSearchQuery, selectedSource),
    enabled: isSearchActive // Sadece arama aktifse bu sorguyu çalıştır
  });

  const topicsList = topics?.topics && Array.isArray(topics.topics) ? topics.topics : [];
  
  // Eğer arama aktifse dönen sonuçları, değilse normal akışı kullan
  const rawArticlesList = isSearchActive 
    ? (searchResults?.articles || []) 
    : (recentArticles?.articles || []);

  const articlesLoading = isSearchActive ? searchLoading : recentLoading;

  // İstemci tarafı (Client-side) Konu (Topic) filtrelemesi
  const filteredArticles = useMemo(() => {
    let list = Array.isArray(rawArticlesList) ? rawArticlesList : [];
    
    // Konu filtresini uygula
    if (selectedTopic !== 'Tümü') {
      const topic = topicsList.find(t => t.name === selectedTopic);
      if (topic) {
        list = list.filter(article => article.topic_id === topic.id);
      }
    }
    
    return list;
  }, [rawArticlesList, topicsList, selectedTopic]);

  return (
    <>
      {/* ========== TOP NAV BAR — Stitch Terminal Monolith ========== */}
      <header className="w-full top-0 sticky bg-background border-b border-outline-variant z-50">
        <nav className="flex justify-between items-center h-16 px-margin-desktop max-w-max-width mx-auto" aria-label="Ana Navigasyon">
          <div className="flex items-center gap-8">
            <span className="text-headline-md font-headline-md text-primary tracking-tighter" aria-label="HaberBot Logo">HaberBot</span>
            <div className="hidden md:flex items-center gap-6" role="menubar">
              <Link to="/" role="menuitem" aria-current="page" className="text-primary font-bold border-b-2 border-primary pb-2 font-body-md text-body-md">Discover</Link>
              <Link to="/" role="menuitem" className="text-on-surface-variant font-medium font-body-md text-body-md hover:text-primary hover:bg-surface-container-low transition-colors duration-150">Feed</Link>
              <Link to="/admin/sources" role="menuitem" className="text-on-surface-variant font-medium font-body-md text-body-md hover:text-primary hover:bg-surface-container-low transition-colors duration-150">Sources</Link>
              <Link to="/admin" role="menuitem" className="text-on-surface-variant font-medium font-body-md text-body-md hover:text-primary hover:bg-surface-container-low transition-colors duration-150">Analytics</Link>
            </div>
          </div>
          <div className="flex items-center gap-4">
            
            {/* Arama Kutusu */}
            <div className="hidden sm:flex items-center bg-[#0d0d0d] border border-[#222222] focus-within:border-primary px-3 py-1.5 rounded-lg gap-2 transition-colors">
              <span className="material-symbols-outlined text-on-surface-variant text-[18px]" aria-hidden="true">search</span>
              <label htmlFor="search-input" className="sr-only">Haberlerde ara</label>
              <input 
                id="search-input"
                className="bg-transparent border-none focus:ring-0 text-body-md text-on-surface p-0 placeholder:text-outline w-64" 
                placeholder="Haberlerde (Türkçe) ara..." 
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              {searchQuery && (
                <button 
                  onClick={() => setSearchQuery('')}
                  className="material-symbols-outlined text-outline hover:text-primary text-[16px]"
                >
                  close
                </button>
              )}
            </div>
            
            <div className="flex items-center gap-2">
              <button className="material-symbols-outlined text-on-surface-variant p-2 hover:bg-surface-container-low transition-colors rounded-full" aria-label="Bildirimler">notifications</button>
              <Link to="/admin" className="material-symbols-outlined text-on-surface-variant p-2 hover:bg-surface-container-low transition-colors rounded-full" aria-label="Yönetici Ayarları">settings</Link>
              <Link to="/login" aria-label="Giriş Yap / Profil" className="h-8 w-8 rounded-full bg-surface-container-high border border-outline-variant flex items-center justify-center overflow-hidden">
                <img 
                  src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80" 
                  alt="Profil" 
                  className="object-cover w-full h-full" 
                />
              </Link>
            </div>
          </div>
        </nav>
      </header>

      {/* ========== MAIN CONTENT ========== */}
      <main className="max-w-max-width mx-auto px-margin-desktop py-8">
        
        {/* Filtreler Bölümü */}
        <section className="mb-12 flex flex-col md:flex-row justify-between items-start md:items-end gap-6" aria-labelledby="summary-heading">
          <div>
            <h1 id="summary-heading" className="font-headline-lg text-headline-lg text-primary mb-6">Günün Özeti</h1>
            
            {/* Konu Filtreleri */}
            <div className="flex flex-col gap-2">
              <span className="text-label-sm text-on-surface-variant uppercase tracking-wider">Konular</span>
              <div role="tablist" aria-label="Konu Kategorileri" className="flex flex-wrap gap-2">
                <button 
                  role="tab"
                  aria-selected={selectedTopic === 'Tümü'}
                  onClick={() => setSelectedTopic('Tümü')}
                  className={`${selectedTopic === 'Tümü' ? 'bg-primary text-on-primary' : 'border border-[#222222] bg-[#0A0A0A] text-on-surface-variant hover:border-[#444444] hover:text-white'} px-4 py-2 rounded-lg font-label-md text-label-md transition-all active:scale-95`}
                >
                  Tümü
                </button>
                {topicsList.map(topic => (
                  <button 
                    key={topic.id}
                    role="tab"
                    aria-selected={selectedTopic === topic.name}
                    onClick={() => setSelectedTopic(topic.name)}
                    className={`${selectedTopic === topic.name ? 'bg-primary text-on-primary' : 'border border-[#222222] bg-[#0A0A0A] text-on-surface-variant hover:border-[#444444] hover:text-white'} px-4 py-2 rounded-lg font-label-md text-label-md transition-all active:scale-95`}
                  >
                    {topic.name}
                  </button>
                ))}
              </div>
            </div>
          </div>

          {/* Kaynak (Source) Filtreleri */}
          <div className="flex flex-col gap-2">
            <span className="text-label-sm text-on-surface-variant uppercase tracking-wider">Kaynaklar</span>
            <div role="tablist" aria-label="Kaynak Filtreleri" className="flex flex-wrap gap-2">
              {['Tümü', 'hackernews', 'rss'].map(source => (
                <button 
                  key={source}
                  role="tab"
                  aria-selected={selectedSource === source}
                  onClick={() => setSelectedSource(source)}
                  className={`${selectedSource === source ? 'bg-[#3b82f6] text-white border-transparent' : 'border border-[#222222] bg-[#0d0d0d] text-on-surface-variant hover:border-[#3b82f6] hover:text-white'} px-4 py-2 rounded-lg font-label-md text-label-md transition-all active:scale-95 capitalize`}
                >
                  {source}
                </button>
              ))}
            </div>
          </div>
        </section>
 
        {/* News Grid — 3-column, matching Stitch */}
        <section 
          id="articles-grid" 
          role="region" 
          aria-live="polite" 
          aria-busy={articlesLoading}
          aria-label="Haber Kartları" 
          className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-gutter"
        >
          {articlesLoading ? (
            Array(6).fill(0).map((_, i) => <ArticleSkeleton key={i} />)
          ) : filteredArticles && filteredArticles.length > 0 ? (
            filteredArticles.map((article, index) => (
              <ArticleCard key={article.id} article={article} index={index} />
            ))
          ) : (
            <div className="col-span-full py-12 flex flex-col items-center justify-center text-center bg-[#0A0A0A] border border-[#222222] rounded-lg p-8">
              <span className="material-symbols-outlined text-[48px] text-[#444444] mb-4">search_off</span>
              <p className="text-on-surface-variant font-body-lg text-body-lg">
                {isSearchActive ? 'Arama kriterlerinize uygun haber bulunamadı.' : 'Bu kategoriye ait haber henüz bulunmamaktadır.'}
              </p>
            </div>
          )}
        </section>

        {/* Load More */}
        {filteredArticles && filteredArticles.length > 0 && !isSearchActive && (
          <div className="mt-12 flex justify-center">
            <button 
              aria-label="Daha fazla haber yükle" 
              className="px-8 py-3 border border-[#222222] hover:border-[#444444] text-primary font-label-md text-label-md transition-all rounded active:scale-95"
            >
              Daha Fazla Haber Yükle
            </button>
          </div>
        )}
      </main>

      {/* ========== FOOTER — Stitch Terminal Monolith ========== */}
      <footer className="w-full py-12 border-t border-outline-variant bg-background mt-12" aria-label="Sayfa Alt Bilgisi">
        <div className="max-w-max-width mx-auto px-margin-desktop flex flex-col md:flex-row justify-between items-center gap-8">
          <div className="flex flex-col gap-2">
            <span className="font-headline-md text-headline-md text-primary">HaberBot</span>
            <p className="text-on-surface-variant font-label-sm text-label-sm">© 2024 HaberBot AI. All rights reserved.</p>
          </div>
          <div className="flex gap-8">
            <a href="#" className="text-on-surface-variant font-label-sm text-label-sm hover:text-primary underline transition-all">Terms</a>
            <a href="#" className="text-on-surface-variant font-label-sm text-label-sm hover:text-primary underline transition-all">Privacy</a>
            <a href="#" className="text-on-surface-variant font-label-sm text-label-sm hover:text-primary underline transition-all">API Docs</a>
            <a href="#" className="text-on-surface-variant font-label-sm text-label-sm hover:text-primary underline transition-all">Status</a>
          </div>
          <div className="flex gap-4">
            <a href="#" className="material-symbols-outlined text-on-surface-variant hover:text-primary transition-colors" aria-label="Terminal">terminal</a>
            <a href="#" className="material-symbols-outlined text-on-surface-variant hover:text-primary transition-colors" aria-label="RSS">rss_feed</a>
            <a href="#" className="material-symbols-outlined text-on-surface-variant hover:text-primary transition-colors" aria-label="Topluluk">group</a>
          </div>
        </div>
      </footer>
    </>
  );
}
