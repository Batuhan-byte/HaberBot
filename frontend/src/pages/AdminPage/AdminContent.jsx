import React, { useState, useEffect, useMemo } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../../services/api';
import AdminLayout from '../../components/Admin/AdminLayout';

export default function AdminContent() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [fetchStatus, setFetchStatus] = useState('');
  const [selectedTopicId, setSelectedTopicId] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [page, setPage] = useState(1);
  const limit = 10;

  const formatDateTime = (dateStr) => {
    if (!dateStr) return '-';
    try {
      const date = new Date(dateStr);
      return date.toLocaleString('tr-TR', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch (e) {
      return '-';
    }
  };

  // Key authorization check
  useEffect(() => {
    const key = localStorage.getItem('admin_api_key');
    if (!key) {
      navigate('/login');
    }
  }, [navigate]);

  // Fetch topics
  const { data: topicsData } = useQuery({
    queryKey: ['topics'],
    queryFn: api.getTopics
  });
  const topicsList = topicsData?.topics && Array.isArray(topicsData.topics) ? topicsData.topics : [];

  // Fetch admin articles (paginated)
  const { data: adminArticlesData, isLoading, refetch } = useQuery({
    queryKey: ['admin_articles', selectedTopicId, page],
    queryFn: () => api.getArticlesAdmin(selectedTopicId, page, limit),
    keepPreviousData: true
  });

  const rawArticles = adminArticlesData?.articles || [];
  const totalArticles = adminArticlesData?.total || 0;
  const totalPages = Math.ceil(totalArticles / limit) || 1;

  // Filter articles locally by search query
  const filteredArticles = useMemo(() => {
    if (!searchQuery) return rawArticles;
    const term = searchQuery.toLowerCase();
    return rawArticles.filter(article => 
      (article.title && article.title.toLowerCase().includes(term)) ||
      (article.title_tr && article.title_tr.toLowerCase().includes(term)) ||
      (article.id && article.id.toLowerCase().includes(term))
    );
  }, [rawArticles, searchQuery]);

  // Mutations for instant visual feedback
  const updateMutation = useMutation({
    mutationFn: ({ id, updates }) => api.updateArticleAdmin(id, updates),
    onSuccess: () => {
      queryClient.invalidateQueries(['admin_articles']);
    }
  });

  const deleteMutation = useMutation({
    mutationFn: (id) => api.deleteArticleAdmin(id),
    onSuccess: () => {
      queryClient.invalidateQueries(['admin_articles']);
    }
  });

  // Action handlers
  const handleToggleApproval = (article) => {
    updateMutation.mutate({
      id: article.id,
      updates: { is_approved: !article.is_approved }
    });
  };

  const handleToggleHiding = (article) => {
    updateMutation.mutate({
      id: article.id,
      updates: { is_hidden: !article.is_hidden }
    });
  };

  const handleDeleteArticle = (article) => {
    if (window.confirm(`"${article.title_tr || article.title}" makalesini kalıcı olarak silmek istediğinizden emin misiniz?`)) {
      deleteMutation.mutate(article.id);
    }
  };

  // Sync pipeline fetcher
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

  const headerActions = (
    <>
      {fetchStatus && <span className="text-blue-400 text-sm font-bold animate-pulse">{fetchStatus}</span>}
      <button onClick={handleFetch} className="bg-blue-600 text-white font-bold text-sm px-6 py-3 rounded-lg flex items-center gap-2 hover:bg-blue-500 hover:shadow-[0_0_20px_rgba(59,130,246,0.3)] transition-all duration-200 active:scale-95 group">
        <span className="material-symbols-outlined group-hover:rotate-180 transition-transform duration-500" data-icon="sync">sync</span>
        Hemen Haber Çek & İşle
      </button>
    </>
  );

  return (
    <AdminLayout
      title="İçerik Yönetim Konsolu"
      subtitle="Haber onay durumlarını kontrol edin, yayından kaldırın veya kalıcı olarak silin."
      actions={headerActions}
    >

          {/* Category Tabs */}
          <div className="flex gap-2 mb-6 overflow-x-auto pb-2 border-b border-zinc-800">
            <button 
              onClick={() => { setSelectedTopicId(''); setPage(1); }} 
              className={`px-4 py-2 text-xs font-extrabold rounded-full border transition-all ${selectedTopicId === '' ? 'bg-blue-600/15 border-blue-500/30 text-blue-400' : 'bg-transparent border-zinc-800 text-zinc-400 hover:border-zinc-700 hover:text-white'}`}
            >
              Tümü
            </button>
            <button 
              onClick={() => { setSelectedTopicId('pending'); setPage(1); }} 
              className={`px-4 py-2 text-xs font-extrabold rounded-full border transition-all whitespace-nowrap ${(selectedTopicId === 'pending' || selectedTopicId.startsWith('pending:')) ? 'bg-amber-500/15 border-amber-500/30 text-amber-400' : 'bg-transparent border-zinc-800 text-zinc-400 hover:border-zinc-700 hover:text-white'}`}
            >
              Onay Bekleyenler
            </button>
            {topicsList.map(topic => (
              <button
                key={topic.id}
                onClick={() => { setSelectedTopicId(topic.id); setPage(1); }}
                className={`px-4 py-2 text-xs font-extrabold rounded-full border transition-all whitespace-nowrap ${selectedTopicId === topic.id ? 'bg-blue-600/15 border-blue-500/30 text-blue-400' : 'bg-transparent border-zinc-800 text-zinc-400 hover:border-zinc-700 hover:text-white'}`}
              >
                {topic.name}
              </button>
            ))}
          </div>

          {/* Table Container */}
          <div className="bg-zinc-950 border border-zinc-800 rounded-xl overflow-hidden flex-grow flex flex-col shadow-lg shadow-black/40">
            <div className="p-4 border-b border-zinc-800 flex justify-between items-center bg-zinc-900/30">
              <div className="flex items-center gap-4">
                <div className="relative group">
                  <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-zinc-500 text-[18px]" data-icon="search">search</span>
                  <input 
                    className="bg-zinc-900 border border-zinc-800 rounded-md pl-9 pr-4 py-1.5 text-xs text-white w-64 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 transition-all outline-none" 
                    placeholder="Haber başlığı veya ID ile ara..." 
                    type="text"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                  />
                </div>

                {/* Secondary Category Filter (only when "Onay Bekleyenler" is active) */}
                {(selectedTopicId === 'pending' || selectedTopicId.startsWith('pending:')) ? (
                  <div className="flex items-center gap-2">
                    <span className="text-[11px] text-zinc-500 font-bold uppercase tracking-wider">Kategori:</span>
                    <select
                      value={selectedTopicId.startsWith('pending:') ? selectedTopicId.split('pending:')[1] : ''}
                      onChange={(e) => {
                        const val = e.target.value;
                        if (val === '') {
                          setSelectedTopicId('pending');
                        } else {
                          setSelectedTopicId(`pending:${val}`);
                        }
                        setPage(1);
                      }}
                      className="bg-zinc-900 border border-zinc-800 rounded-md px-3 py-1.5 text-xs text-white outline-none focus:ring-1 focus:ring-amber-500 focus:border-amber-500 font-bold cursor-pointer transition-all hover:border-zinc-700"
                    >
                      <option value="">Tümü</option>
                      {topicsList.map(topic => (
                        <option key={topic.id} value={topic.id}>{topic.name}</option>
                      ))}
                    </select>
                  </div>
                ) : null}
              </div>
              <div className="flex items-center gap-2 text-zinc-400 text-xs font-bold">
                <span>Tabloda Gösterilen: <strong className="text-white">{filteredArticles.length}</strong> / Sistemde Kayıtlı: <strong className="text-white">{totalArticles}</strong></span>
              </div>
            </div>

            {/* Table */}
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="bg-zinc-900/60 border-b border-zinc-800 text-zinc-500">
                    <th className="py-4 px-6 text-xs font-extrabold uppercase tracking-wider">Haber Başlığı</th>
                    <th className="py-4 px-6 text-xs font-extrabold uppercase tracking-wider">Kategori & Kaynak</th>
                    <th className="py-4 px-6 text-xs font-extrabold uppercase tracking-wider">Durumlar</th>
                    <th className="py-4 px-6 text-xs font-extrabold uppercase tracking-wider text-right">İşlemler</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-zinc-900">
                  {isLoading ? (
                    <tr>
                      <td colSpan="4" className="py-12 px-6 text-center text-zinc-500">
                        <div className="flex items-center justify-center gap-2">
                          <span className="material-symbols-outlined animate-spin text-blue-500" data-icon="sync">sync</span>
                          Haber verileri yükleniyor...
                        </div>
                      </td>
                    </tr>
                  ) : filteredArticles.length > 0 ? (
                    filteredArticles.map(article => (
                      <tr key={article.id} className="hover:bg-zinc-900/40 transition-colors duration-150 group">
                        {/* Title and ID */}
                        <td className="py-4 px-6 max-w-[460px]">
                          <div className="flex flex-col gap-1.5">
                            <Link 
                              to={`/haber/${article.id}`} 
                              target="_blank" 
                              rel="noopener noreferrer" 
                              className="text-sm text-white font-bold block truncate hover:text-blue-400 transition-colors hover:underline"
                            >
                              {article.title_tr || article.title}
                            </Link>
                            {article.title_tr && (
                              <span className="text-[11px] text-zinc-500 block truncate font-medium">Orj: {article.title}</span>
                            )}
                            
                            <div className="flex flex-wrap items-center gap-x-3 gap-y-1 text-[11px] text-zinc-500 font-medium">
                              <span className="text-[10px] text-zinc-600 font-mono select-all">ID: {article.id}</span>
                              <span className="flex items-center gap-1">
                                <span className="material-symbols-outlined text-[12px] text-zinc-500" style={{ verticalAlign: 'middle' }}>calendar_today</span>
                                Giriş: <strong className="text-zinc-400">{formatDateTime(article.created_at)}</strong>
                              </span>
                              {article.is_approved && article.approved_at ? (
                                <span className="flex items-center gap-1 text-emerald-400 bg-emerald-950/20 px-1.5 py-0.5 rounded border border-emerald-500/10">
                                  <span className="material-symbols-outlined text-[12px]" style={{ verticalAlign: 'middle' }}>check_circle</span>
                                  Onay: <strong>{formatDateTime(article.approved_at)}</strong>
                                </span>
                              ) : (
                                <span className="flex items-center gap-1 text-amber-400 bg-amber-950/20 px-1.5 py-0.5 rounded border border-amber-500/10">
                                  <span className="material-symbols-outlined text-[12px]" style={{ verticalAlign: 'middle' }}>pending</span>
                                  Onay: <strong>Bekliyor</strong>
                                </span>
                              )}
                            </div>
                          </div>
                        </td>

                        {/* Category & Source */}
                        <td className="py-4 px-6">
                          <div className="flex flex-wrap gap-1.5 items-center">
                            <span className="px-2 py-0.5 bg-blue-600/10 border border-blue-500/20 text-blue-400 rounded text-[10px] font-extrabold uppercase tracking-wide">
                              {topicsList.find(t => t.id === article.topic_id)?.name || 'Teknoloji'}
                            </span>
                            <span className="px-2 py-0.5 bg-zinc-800 border border-zinc-700 text-zinc-300 rounded text-[10px] font-extrabold uppercase tracking-wide">
                              {article.source === 'hackernews' ? 'HN' : 'RSS'}
                            </span>
                          </div>
                        </td>

                        {/* Status indicators */}
                        <td className="py-4 px-6">
                          <div className="flex flex-wrap gap-1.5">


                            {/* Approval Status */}
                            {article.is_approved ? (
                              <span className="px-2 py-0.5 bg-green-500/15 border border-green-500/30 text-green-400 rounded text-[10px] font-extrabold uppercase tracking-wide">Onaylı</span>
                            ) : (
                              <span className="px-2 py-0.5 bg-amber-500/15 border border-amber-500/30 text-amber-400 rounded text-[10px] font-extrabold uppercase tracking-wide">Onay Bekliyor</span>
                            )}

                            {/* Visibility Status */}
                            {article.is_hidden ? (
                              <span className="px-2 py-0.5 bg-red-500/15 border border-red-500/30 text-red-400 rounded text-[10px] font-extrabold uppercase tracking-wide">Gizli</span>
                            ) : (
                              <span className="px-2 py-0.5 bg-blue-500/15 border border-blue-500/30 text-blue-400 rounded text-[10px] font-extrabold uppercase tracking-wide">Yayında</span>
                            )}
                          </div>
                        </td>

                        {/* Action buttons */}
                        <td className="py-4 px-6 text-right">
                          <div className="flex gap-1.5 justify-end items-center opacity-70 group-hover:opacity-100 transition-opacity">
                            {/* Approve button */}
                            <button 
                              onClick={() => handleToggleApproval(article)}
                              className={`p-1.5 rounded transition-colors ${article.is_approved ? 'bg-green-500/10 text-green-400 hover:bg-green-500/20' : 'bg-zinc-900 text-zinc-500 hover:bg-zinc-800 hover:text-white'}`}
                              title={article.is_approved ? 'Onayı Kaldır' : 'Onayla'}
                              disabled={updateMutation.isPending}
                            >
                              <span className="material-symbols-outlined text-[18px]">
                                {article.is_approved ? 'check_circle' : 'unpublished'}
                              </span>
                            </button>

                            {/* Hide button */}
                            <button 
                              onClick={() => handleToggleHiding(article)}
                              className={`p-1.5 rounded transition-colors ${article.is_hidden ? 'bg-red-500/10 text-red-400 hover:bg-red-500/20' : 'bg-blue-500/10 text-blue-400 hover:bg-blue-500/20'}`}
                              title={article.is_hidden ? 'Yayına Al' : 'Gizle'}
                              disabled={updateMutation.isPending}
                            >
                              <span className="material-symbols-outlined text-[18px]">
                                {article.is_hidden ? 'visibility_off' : 'visibility'}
                              </span>
                            </button>

                            {/* Delete button */}
                            <button 
                              onClick={() => handleDeleteArticle(article)}
                              className="p-1.5 bg-zinc-900 text-zinc-500 hover:bg-red-500/20 hover:text-red-400 rounded transition-colors"
                              title="Kalıcı Olarak Sil"
                              disabled={deleteMutation.isPending}
                            >
                              <span className="material-symbols-outlined text-[18px]">delete</span>
                            </button>
                          </div>
                        </td>
                      </tr>
                    ))
                  ) : (
                    <tr>
                      <td colSpan="4" className="py-12 px-6 text-center text-zinc-500">
                        {searchQuery ? 'Aramanızla eşleşen makale bulunamadı.' : 'Bu kategoriye ait makale bulunmuyor.'}
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>

            {/* Pagination Controls */}
            {totalPages > 1 && (
              <div className="p-4 border-t border-zinc-800 bg-zinc-900/30 flex justify-between items-center mt-auto">
                <span className="text-zinc-500 text-xs font-bold">Sayfa <span className="text-white">{page}</span> / <span className="text-white">{totalPages}</span></span>
                <div className="flex gap-2">
                  <button 
                    disabled={page === 1}
                    onClick={() => setPage(page - 1)}
                    className="px-4 py-1.5 border border-zinc-800 hover:border-zinc-700 disabled:opacity-40 disabled:hover:border-zinc-800 rounded transition-all text-xs font-extrabold text-zinc-300"
                  >
                    Önceki
                  </button>
                  <button 
                    disabled={page === totalPages}
                    onClick={() => setPage(page + 1)}
                    className="px-4 py-1.5 border border-zinc-800 hover:border-zinc-700 disabled:opacity-40 disabled:hover:border-zinc-800 rounded transition-all text-xs font-extrabold text-zinc-300"
                  >
                    Sonraki
                  </button>
                </div>
              </div>
            )}
          </div>
    </AdminLayout>
  );
}
