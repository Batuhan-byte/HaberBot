import React, { useState, useEffect } from 'react';
import { Link, useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../../services/api';

function ArticlePageSkeleton() {
  return (
    <div className="max-w-max-width mx-auto px-margin-desktop py-12 animate-pulse" aria-hidden="true">
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-gutter">
        <div className="lg:col-span-8">
          <div className="flex items-center gap-3 mb-6">
            <div className="h-6 w-16 bg-neutral-800 rounded"></div>
            <div className="h-4 w-28 bg-neutral-800 rounded"></div>
          </div>
          <div className="h-12 w-full bg-neutral-800 rounded mb-6"></div>
          <div className="h-[380px] w-full bg-neutral-800 rounded-xl mb-10"></div>
          <div className="space-y-4">
            <div className="h-4 w-full bg-neutral-800 rounded"></div>
            <div className="h-4 w-full bg-neutral-800 rounded"></div>
            <div className="h-4 w-5/6 bg-neutral-800 rounded"></div>
          </div>
        </div>
        <div className="lg:col-span-4">
          <div className="h-48 w-full bg-neutral-800 rounded"></div>
        </div>
      </div>
    </div>
  );
}

export default function ArticlePage() {
  const { id } = useParams();

  const { data: article, isLoading, error } = useQuery({
    queryKey: ['article', id],
    queryFn: () => api.getArticle(id),
  });

  const queryClient = useQueryClient();
  const summaryMutation = useMutation({
    mutationFn: (articleId) => api.generateSummary(articleId),
    onSuccess: (data) => {
      // Update cache with the newly generated summary
      queryClient.setQueryData(['article', id], (oldData) => ({
        ...oldData,
        summary_tr: data.summary,
      }));
    }
  });

  if (isLoading) {
    return (
      <>
        <header className="w-full top-0 sticky bg-[#050505] border-b border-white/5 z-50">
          <nav className="flex justify-between items-center h-16 px-margin-desktop max-w-max-width mx-auto">
            <span className="text-white font-extrabold text-2xl tracking-tighter">HaberBot</span>
          </nav>
        </header>
        <ArticlePageSkeleton />
      </>
    );
  }

  if (error || !article) {
    return (
      <div className="text-center py-20 text-error flex flex-col items-center justify-center gap-4">
        <span className="material-symbols-outlined text-[48px]">error</span>
        <p className="font-body-lg text-body-lg text-on-surface-variant">Haber bulunamadı veya bir hata oluştu.</p>
        <Link to="/" className="px-6 py-2 bg-primary text-on-primary rounded font-label-md text-label-md active:scale-95 transition-all">
          Ana Sayfaya Dön
        </Link>
      </div>
    );
  }

  return (
    <>
      <header className="w-full top-0 sticky bg-[#050505] border-b border-white/5 z-50">
        <nav className="flex justify-between items-center h-16 px-margin-desktop max-w-max-width mx-auto" aria-label="Ana Navigasyon">
          <div className="flex items-center gap-8">
            <span className="text-white font-extrabold text-2xl tracking-tighter" aria-label="HaberBot Logo">HaberBot</span>
            <div className="hidden md:flex items-center gap-6" role="menubar">
              <Link to="/" role="menuitem" className="text-neutral-400 hover:text-white font-medium text-sm transition-colors duration-150">Discover</Link>
              <Link to="/" role="menuitem" aria-current="page" className="text-white font-bold border-b-2 border-white pb-2 text-sm">Feed</Link>
              <Link to="/admin/sources" role="menuitem" className="text-neutral-400 hover:text-white font-medium text-sm transition-colors duration-150">Sources</Link>
              <Link to="/admin" role="menuitem" className="text-neutral-400 hover:text-white font-medium text-sm transition-colors duration-150">Analytics</Link>
            </div>
          </div>
          <div className="flex items-center gap-4">
            <div className="hidden sm:flex items-center bg-[#121214] border border-white/10 px-3 py-1.5 rounded-lg gap-2 w-64">
              <span className="material-symbols-outlined text-neutral-400 text-[18px]" aria-hidden="true">search</span>
              <label htmlFor="article-search-input" className="sr-only">Haberlerde ara</label>
              <input 
                id="article-search-input"
                className="bg-transparent border-none focus:ring-0 text-sm text-white p-0 placeholder:text-neutral-500" 
                placeholder="Haberlerde ara..." 
                type="text"
              />
            </div>
            <div className="flex items-center gap-3">
              <span className="material-symbols-outlined text-neutral-400 hover:text-white transition-colors cursor-pointer text-[22px]" aria-hidden="true">notifications</span>
              <Link to="/admin" aria-label="Yönetici Ayarları" className="material-symbols-outlined text-neutral-400 hover:text-white p-1 transition-colors text-[22px]">settings</Link>
              <Link to="/login" aria-label="Giriş Yap / Profil" className="h-8 w-8 rounded-full overflow-hidden border border-white/15 flex items-center justify-center">
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

      <main className="max-w-max-width mx-auto px-margin-desktop py-12">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-gutter">
          
          <article className="lg:col-span-8" aria-labelledby="article-title">
            <header className="mb-10">
              <div className="flex items-center gap-3 mb-6">
                <span className="px-2 py-0.5 bg-surface-container-highest border border-outline-variant font-label-sm text-label-sm rounded uppercase tracking-wider text-primary">{article.source || 'Haber'}</span>
                <span className="text-on-surface-variant font-label-md text-label-md">{article.created_at ? new Date(article.created_at).toLocaleString('tr-TR') : new Date().toLocaleString('tr-TR')}</span>
              </div>
              <h1 id="article-title" className="font-headline-lg text-headline-lg text-primary mb-6 leading-tight">
                {article.title_tr || article.title}
              </h1>
              
              {!article.summary_tr && (
                <button
                  onClick={() => summaryMutation.mutate(id)}
                  disabled={summaryMutation.isPending}
                  className="mb-6 px-6 py-3 bg-[#1A1A1C] hover:bg-[#252528] active:scale-95 transition-all text-primary font-bold rounded-full border border-primary/20 hover:border-primary/50 flex items-center gap-3 shadow-[0_0_15px_rgba(255,255,255,0.05)]"
                >
                  <span className={`material-symbols-outlined ${summaryMutation.isPending ? 'animate-spin' : ''}`}>
                    {summaryMutation.isPending ? 'sync' : 'auto_awesome'}
                  </span>
                  <span>{summaryMutation.isPending ? 'Yapay Zeka Özetliyor...' : '✨ Yapay Zeka ile Özetle'}</span>
                </button>
              )}
              <div className="flex items-center justify-between border-y border-outline-variant py-4">
                <div className="flex items-center gap-4">
                  <div className="w-10 h-10 rounded bg-surface-container border border-outline-variant flex items-center justify-center">
                    <span className="material-symbols-outlined text-primary" aria-hidden="true">memory</span>
                  </div>
                  <div>
                    <div className="font-label-md text-label-md text-primary">{article.source || 'Bilinmeyen Kaynak'}</div>
                    <div className="font-label-sm text-label-sm text-on-surface-variant">Yayıncı: HaberBot AI</div>
                  </div>
                </div>
                <div className="flex gap-2">
                  <a 
                    href={article.original_url || article.url} 
                    target="_blank" 
                    rel="noopener noreferrer" 
                    className="p-2 border border-outline-variant hover:border-white transition-colors rounded flex items-center justify-center"
                    aria-label="Orijinal kaynağa harici sekmede git"
                  >
                    <span className="material-symbols-outlined text-[20px]" aria-hidden="true">open_in_new</span>
                  </a>
                </div>
              </div>
            </header>

            {/* Feature Image */}
            <div className="mb-10 rounded-xl overflow-hidden border border-outline-variant shadow-lg group relative h-[380px] w-full bg-[#111111]">
              <img 
                src={article.image_url || "https://images.unsplash.com/photo-1488590528505-98d2b5aba04b?auto=format&fit=crop&w=1200&h=600&q=80"} 
                alt={article.title_tr || article.title || "Haber görseli"} 
                className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105" 
              />
              {!article.image_url && (
                <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/30 to-transparent flex items-end p-6">
                  <span className="text-sm font-label-sm px-3 py-1 bg-primary/20 backdrop-blur-md text-primary border border-primary/30 rounded-full font-medium">
                    Haber Gelişmesi
                  </span>
                </div>
              )}
            </div>
 
            <section className="news-content" aria-label="Haber İçeriği">
              {(article.summary_tr || article.summary) && (
                <div className="mb-8 p-6 bg-surface-container-low border border-outline-variant rounded-xl" role="region" aria-label="Yapay Zeka Özeti">
                  <div className="flex items-center gap-2 mb-4 text-on-surface-variant font-label-md text-label-md">
                    <span className="material-symbols-outlined text-primary scale-75" aria-hidden="true">auto_awesome</span>
                    <span className="text-primary font-semibold tracking-wide">YAPAY ZEKA ÖZETİ</span>
                  </div>
                  <div className="font-body-lg text-body-lg leading-[1.8] text-on-surface animate-fade-in">
                    <p className="font-semibold text-white/90">
                      {article.summary_tr || article.summary}
                    </p>
                  </div>
                </div>
              )}
              
              {summaryMutation.isPending && !article.summary_tr && (
                <div className="mb-8 p-6 bg-surface-container-low border border-outline-variant rounded-xl animate-pulse">
                  <div className="flex items-center gap-2 mb-4">
                     <div className="w-5 h-5 rounded-full bg-primary/30"></div>
                     <div className="h-4 w-32 bg-primary/20 rounded"></div>
                  </div>
                  <div className="space-y-3">
                    <div className="h-4 w-full bg-neutral-800 rounded"></div>
                    <div className="h-4 w-full bg-neutral-800 rounded"></div>
                    <div className="h-4 w-3/4 bg-neutral-800 rounded"></div>
                  </div>
                </div>
              )}
 
              <div 
                className="font-body-lg text-body-lg leading-[1.8] text-on-surface space-y-6"
                dangerouslySetInnerHTML={{ 
                  __html: article.content_tr || article.original_content || article.content || "<p>Haberin detayı henüz çekilmedi veya bulunmuyor.</p>" 
                }}
              />
            </section>

          </article>

          <aside className="lg:col-span-4 space-y-8" aria-label="Yan Menü Bilgileri">
            <div className="bg-surface-container-low border border-outline-variant p-6 rounded-xl">
              <h3 className="font-label-md text-label-md text-on-surface-variant uppercase mb-4">Haber Kaynağı</h3>
              <a 
                className="inline-flex items-center justify-between w-full bg-primary text-on-primary px-4 py-3 rounded-lg font-label-md text-label-md hover:opacity-90 active:scale-95 transition-all shadow-lg shadow-primary/20" 
                href={article.url || article.original_url} 
                target="_blank" 
                rel="noopener noreferrer"
                aria-label="Haberin orijinal kaynağına git"
              >
                <span>Kaynağa Git</span>
                <span className="material-symbols-outlined text-[18px]" aria-hidden="true">open_in_new</span>
              </a>
              <div className="mt-4 flex items-center gap-2 text-on-surface-variant font-label-sm text-label-sm">
                <span className="material-symbols-outlined text-[14px]" aria-hidden="true">link</span>
                <span className="truncate" title={article.url || article.original_url}>{article.url || article.original_url}</span>
              </div>
            </div>
          </aside>
        </div>
      </main>

      <footer className="w-full py-12 border-t border-white/5 bg-[#050505] mt-12" aria-label="Sayfa Alt Bilgisi">
        <div className="max-w-max-width mx-auto px-margin-desktop flex flex-col md:flex-row justify-between items-center gap-8">
          <div className="flex flex-col gap-2">
            <span className="text-white font-extrabold text-2xl tracking-tighter">HaberBot</span>
            <p className="text-neutral-500 font-label-sm text-label-sm">© 2024 HaberBot AI. All rights reserved.</p>
          </div>
          <div className="flex flex-col md:flex-row items-center gap-8">
            <div className="flex gap-6">
              <a href="#" className="text-neutral-400 hover:text-white text-sm transition-colors">Terms</a>
              <a href="#" className="text-neutral-400 hover:text-white text-sm transition-colors">Privacy</a>
              <a href="#" className="text-neutral-400 hover:text-white text-sm transition-colors">API Docs</a>
              <a href="#" className="text-neutral-400 hover:text-white text-sm transition-colors">Status</a>
            </div>
            <div className="flex gap-4 text-neutral-400">
              <span className="material-symbols-outlined hover:text-white cursor-pointer transition-colors text-[20px]">feed</span>
              <span className="material-symbols-outlined hover:text-white cursor-pointer transition-colors text-[20px]">rss_feed</span>
              <span className="material-symbols-outlined hover:text-white cursor-pointer transition-colors text-[20px]">group</span>
            </div>
          </div>
        </div>
      </footer>
    </>
  );
}
