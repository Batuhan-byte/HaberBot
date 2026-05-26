import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { api } from '../../services/api';

export default function AdminContent() {
  const navigate = useNavigate();
  const [fetchStatus, setFetchStatus] = useState('');

  useEffect(() => {
    const key = localStorage.getItem('admin_api_key');
    if (!key) {
      navigate('/login');
    }
  }, [navigate]);

  const { data: articlesData, isLoading, refetch } = useQuery({
    queryKey: ['admin_articles'],
    queryFn: api.getArticles
  });

  const articles = articlesData?.articles && Array.isArray(articlesData.articles) ? articlesData.articles : [];

  const handleFetch = async () => {
    try {
      setFetchStatus('Haberler çekiliyor...');
      await api.triggerFetch();
      setFetchStatus('Haberler çekildi. AI özetleme başlatılıyor...');
      await api.triggerProcess();
      setFetchStatus('İşlem tamamlandı.');
      refetch();
      setTimeout(() => setFetchStatus(''), 3000);
    } catch (error) {
      console.error(error);
      setFetchStatus('Hata oluştu!');
      setTimeout(() => setFetchStatus(''), 3000);
    }
  };

  return (
    <>
      <div className="flex min-h-screen">
        <aside className="h-screen w-64 fixed left-0 top-0 bg-surface-container-lowest border-r border-outline-variant flex flex-col py-6 px-4 gap-4 z-50">
          <div className="flex flex-col gap-1 mb-6 px-2">
            <h1 className="font-headline-md text-headline-md text-primary">Admin Console</h1>
            <p className="font-label-sm text-label-sm text-on-surface-variant opacity-60">v1.0.4</p>
          </div>
          <nav className="flex flex-col gap-1 flex-grow">
            <Link to="/admin" className="flex items-center gap-3 py-2.5 px-3 rounded text-on-surface-variant hover:bg-surface-container hover:text-primary transition-all duration-150 cursor-pointer active:opacity-80 font-label-md text-label-md">
              <span className="material-symbols-outlined" data-icon="dashboard">dashboard</span>
              Dashboard
            </Link>
            <Link to="/admin/sources" className="flex items-center gap-3 py-2.5 px-3 rounded text-on-surface-variant hover:bg-surface-container hover:text-primary transition-all duration-150 cursor-pointer active:opacity-80 font-label-md text-label-md">
              <span className="material-symbols-outlined" data-icon="rss_feed">rss_feed</span>
              RSS Sources
            </Link>
            <Link to="/admin/content" className="flex items-center gap-3 py-2.5 px-3 rounded bg-surface-container-high text-primary border-r-2 border-primary transition-all duration-150 cursor-pointer active:opacity-80 font-label-md text-label-md">
              <span className="material-symbols-outlined" data-icon="article">article</span>
              Content
            </Link>
            <Link to="/" className="flex items-center gap-3 py-2.5 px-3 rounded text-on-surface-variant hover:bg-surface-container hover:text-primary transition-all duration-150 cursor-pointer active:opacity-80 font-label-md text-label-md">
              <span className="material-symbols-outlined" data-icon="home">home</span>
              Siteye Dön
            </Link>
          </nav>
        </aside>

        <main className="flex-grow ml-64 p-margin-desktop min-h-screen flex flex-col">
          <header className="flex justify-between items-end mb-8">
            <div>
              <h2 className="font-headline-lg text-headline-lg text-primary tracking-tight">İçerik Yönetimi</h2>
              <p className="text-on-surface-variant mt-1 font-body-md opacity-80">Haber akışlarını denetleyin ve AI özet kalitesini izleyin.</p>
            </div>
            <div className="flex items-center gap-4">
              {fetchStatus && <span className="text-primary font-label-md">{fetchStatus}</span>}
              <button onClick={handleFetch} className="bg-primary text-on-primary font-label-md text-label-md px-6 py-3 rounded-lg flex items-center gap-2 hover:bg-white hover:shadow-[0_0_20px_rgba(255,255,255,0.2)] transition-all duration-200 active:scale-95 group">
                <span className="material-symbols-outlined group-hover:rotate-180 transition-transform duration-500" data-icon="sync">sync</span>
                Hemen Haber Çek & İşle
              </button>
            </div>
          </header>

          <div className="bg-surface-container-lowest border border-outline-variant rounded-xl overflow-hidden flex-grow flex flex-col">
            <div className="p-4 border-b border-outline-variant flex justify-between items-center bg-surface-container-low">
              <div className="flex items-center gap-4">
                <div className="relative group">
                  <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[20px]" data-icon="search">search</span>
                  <input className="bg-background border border-outline-variant rounded-md pl-10 pr-4 py-1.5 font-body-md text-on-surface w-64 focus:ring-1 focus:ring-primary focus:border-primary transition-all outline-none" placeholder="Başlıklarda ara..." type="text"/>
                </div>
              </div>
              <div className="flex items-center gap-2 text-on-surface-variant font-label-sm text-label-sm">
                <span>Gösteriliyor: <strong>{articles?.length || 0} Haber</strong></span>
              </div>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="bg-surface-container-low border-b border-outline-variant">
                    <th className="py-4 px-6 font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Haber Başlığı</th>
                    <th className="py-4 px-6 font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Kaynak</th>
                    <th className="py-4 px-6 font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Durum</th>
                    <th className="py-4 px-6 font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Yayın Tarihi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-[#111111]">
                  {isLoading ? (
                    <tr><td colSpan="4" className="py-5 px-6 text-center text-on-surface-variant">Yükleniyor...</td></tr>
                  ) : articles?.length > 0 ? (
                    articles.map(article => (
                      <tr key={article.id} className="hover:bg-surface-container-low transition-colors duration-150 group">
                        <td className="py-5 px-6">
                          <div className="flex flex-col gap-1">
                            <Link to={`/haber/${article.id}`} className="font-body-md text-on-surface font-medium block truncate max-w-[400px] hover:underline">
                              {article.title}
                            </Link>
                            <span className="text-[11px] text-on-surface-variant/50 font-label-sm">ID: {article.id}</span>
                          </div>
                        </td>
                        <td className="py-5 px-6">
                          <span className="font-label-md text-label-md px-2 py-1 bg-surface-container-high border border-outline-variant rounded">{article.source || 'Bilinmiyor'}</span>
                        </td>
                        <td className="py-5 px-6">
                          {article.summary_tr ? (
                            <span className="font-label-md text-label-md text-primary">Özetlendi</span>
                          ) : (
                            <span className="font-label-md text-label-md text-error">Özet Bekliyor</span>
                          )}
                        </td>
                        <td className="py-5 px-6 text-on-surface-variant font-label-md text-label-md">
                          {new Date(article.created_at || article.fetched_at).toLocaleString()}
                        </td>
                      </tr>
                    ))
                  ) : (
                    <tr><td colSpan="4" className="py-5 px-6 text-center text-on-surface-variant">Kayıt bulunamadı.</td></tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>
        </main>
      </div>
    </>
  );
}
