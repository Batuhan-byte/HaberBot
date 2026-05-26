import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../../services/api';

export default function AdminSources() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [isAdding, setIsAdding] = useState(false);
  const [newTopic, setNewTopic] = useState({ name: '', slug: '' });

  useEffect(() => {
    const key = localStorage.getItem('admin_api_key');
    if (!key) {
      navigate('/login');
    }
  }, [navigate]);

  const { data: topicsData, isLoading } = useQuery({
    queryKey: ['admin_topics'],
    queryFn: api.getTopics
  });

  const topics = topicsData?.topics && Array.isArray(topicsData.topics) ? topicsData.topics : [];

  const addMutation = useMutation({
    mutationFn: api.createTopic,
    onSuccess: () => {
      queryClient.invalidateQueries(['admin_topics']);
      setIsAdding(false);
      setNewTopic({ name: '', slug: '' });
    }
  });

  const deleteMutation = useMutation({
    mutationFn: api.deleteTopic,
    onSuccess: () => {
      queryClient.invalidateQueries(['admin_topics']);
    }
  });

  const handleAdd = (e) => {
    e.preventDefault();
    if (newTopic.name && newTopic.slug) {
      addMutation.mutate(newTopic);
    }
  };

  const handleDelete = (id) => {
    if (window.confirm('Bu kaynağı silmek istediğinize emin misiniz?')) {
      deleteMutation.mutate(id);
    }
  };

  return (
    <>
      <aside className="h-screen w-64 fixed left-0 top-0 bg-surface-container-lowest border-r border-outline-variant flex flex-col py-6 px-4 gap-4 z-50">
        <div className="flex items-center gap-3 px-2 mb-4">
          <div className="w-10 h-10 bg-primary flex items-center justify-center rounded">
            <span className="material-symbols-outlined text-background">robot_2</span>
          </div>
          <div>
            <h1 className="font-headline-md text-headline-md text-primary leading-none">HaberBot</h1>
            <p className="font-label-sm text-label-sm text-on-surface-variant opacity-60">v1.0.4</p>
          </div>
        </div>
        <nav className="flex-1 flex flex-col gap-1">
          <Link to="/admin" className="flex items-center gap-3 py-2 px-3 text-on-surface-variant hover:bg-surface-container hover:text-primary transition-all duration-150 rounded cursor-pointer active:opacity-80">
            <span className="material-symbols-outlined" data-icon="dashboard">dashboard</span>
            <span className="font-label-md text-label-md">Dashboard</span>
          </Link>
          <Link to="/admin/sources" className="flex items-center gap-3 py-2 px-3 bg-surface-container-high text-primary border-r-2 border-primary rounded cursor-pointer active:opacity-80">
            <span className="material-symbols-outlined" data-icon="rss_feed">rss_feed</span>
            <span className="font-label-md text-label-md">RSS Sources</span>
          </Link>
          <Link to="/admin/content" className="flex items-center gap-3 py-2 px-3 text-on-surface-variant hover:bg-surface-container hover:text-primary transition-all duration-150 rounded cursor-pointer active:opacity-80">
            <span className="material-symbols-outlined" data-icon="article">article</span>
            <span className="font-label-md text-label-md">Content</span>
          </Link>
          <Link to="/" className="flex items-center gap-3 py-2 px-3 text-on-surface-variant hover:bg-surface-container hover:text-primary transition-all duration-150 rounded cursor-pointer active:opacity-80">
            <span className="material-symbols-outlined" data-icon="home">home</span>
            <span className="font-label-md text-label-md">Siteye Dön</span>
          </Link>
        </nav>
        <button onClick={() => setIsAdding(true)} className="mt-4 bg-primary text-on-primary py-2 px-4 rounded font-label-md text-label-md flex items-center justify-center gap-2 hover:opacity-90 active:scale-95 transition-all">
          <span className="material-symbols-outlined text-[18px]">add</span>
          Add New Source
        </button>
      </aside>

      <main className="ml-64 min-h-screen bg-background p-margin-desktop flex flex-col">
        <header className="flex justify-between items-end mb-8">
          <div>
            <h2 className="font-headline-lg text-headline-lg text-primary tracking-tight">Source Manager</h2>
            <p className="text-on-surface-variant mt-1">Configure and monitor automated data ingestion streams.</p>
          </div>
        </header>

        {isAdding && (
          <div className="mb-8 bg-surface-container-lowest border border-outline-variant p-6 rounded-lg">
            <h3 className="font-headline-md text-primary mb-4">Yeni Kaynak Ekle</h3>
            <form onSubmit={handleAdd} className="flex gap-4 items-end">
              <div>
                <label className="block text-label-sm text-on-surface-variant mb-1">Adı</label>
                <input 
                  type="text" 
                  value={newTopic.name} 
                  onChange={e => setNewTopic({...newTopic, name: e.target.value})}
                  className="bg-background border border-outline-variant rounded px-3 py-2 text-on-surface w-64 focus:border-primary outline-none" 
                  placeholder="Örn: HackerNews"
                  required
                />
              </div>
              <div>
                <label className="block text-label-sm text-on-surface-variant mb-1">Slug (Benzersiz Tanımlayıcı)</label>
                <input 
                  type="text" 
                  value={newTopic.slug} 
                  onChange={e => setNewTopic({...newTopic, slug: e.target.value})}
                  className="bg-background border border-outline-variant rounded px-3 py-2 text-on-surface w-64 focus:border-primary outline-none" 
                  placeholder="Örn: hackernews"
                  required
                />
              </div>
              <button type="submit" disabled={addMutation.isPending} className="bg-primary text-on-primary px-6 py-2.5 rounded font-label-md hover:opacity-90 transition-all">
                {addMutation.isPending ? 'Ekleniyor...' : 'Kaydet'}
              </button>
              <button type="button" onClick={() => setIsAdding(false)} className="border border-outline-variant text-on-surface-variant px-6 py-2.5 rounded font-label-md hover:bg-surface-container transition-all">
                İptal
              </button>
            </form>
          </div>
        )}

        <div className="bg-surface-container-lowest border border-outline-variant flex flex-col overflow-hidden">
          <div className="overflow-x-auto custom-scrollbar">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="border-b border-outline-variant">
                  <th className="px-6 py-4 font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Source Name</th>
                  <th className="px-6 py-4 font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Slug</th>
                  <th className="px-6 py-4 font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-surface-container-high">
                {isLoading ? (
                  <tr><td colSpan="3" className="px-6 py-4 text-center text-on-surface-variant">Yükleniyor...</td></tr>
                ) : topics?.length > 0 ? (
                  topics.map(topic => (
                    <tr key={topic.id} className="hover:bg-surface-container-low transition-colors group">
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-3">
                          <div className="w-2 h-2 rounded-full bg-emerald-500"></div>
                          <span className="font-label-md text-label-md text-primary">{topic.name}</span>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <span className="font-label-sm text-label-sm text-on-surface-variant font-mono opacity-80">{topic.slug}</span>
                      </td>
                      <td className="px-6 py-4 text-right">
                        <div className="flex justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                          <button onClick={() => handleDelete(topic.id)} disabled={deleteMutation.isPending} className="p-1 hover:text-error text-on-surface-variant transition-colors" title="Delete">
                            <span className="material-symbols-outlined text-[20px]">delete</span>
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr><td colSpan="3" className="px-6 py-4 text-center text-on-surface-variant">Kaynak bulunamadı.</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </main>
    </>
  );
}
