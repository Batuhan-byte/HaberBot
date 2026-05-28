import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../../services/api';
import AdminLayout from '../../components/Admin/AdminLayout';

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

  const headerActions = (
    <button
      onClick={() => setIsAdding((prev) => !prev)}
      className="bg-blue-600 text-white font-bold text-sm px-6 py-3 rounded-lg flex items-center gap-2 hover:bg-blue-500 hover:shadow-[0_0_20px_rgba(59,130,246,0.3)] transition-all duration-200 active:scale-95 group"
    >
      <span className="material-symbols-outlined text-[18px]" data-icon="add">add</span>
      {isAdding ? 'Formu Kapat' : 'Yeni Kaynak Ekle'}
    </button>
  );

  return (
    <AdminLayout
      title="Source Manager"
      subtitle="Configure and monitor automated data ingestion streams."
      actions={headerActions}
    >

      {isAdding && (
        <div className="mb-8 bg-zinc-950 border border-zinc-800 rounded-xl p-6 shadow-lg shadow-black/40">
          <h3 className="text-lg font-bold text-white mb-4">Yeni Kaynak Ekle</h3>
          <form onSubmit={handleAdd} className="flex flex-wrap gap-4 items-end">
            <div>
              <label className="block text-xs text-zinc-500 font-bold uppercase tracking-wider mb-1">Adı</label>
              <input 
                type="text" 
                value={newTopic.name} 
                onChange={e => setNewTopic({...newTopic, name: e.target.value})}
                className="bg-zinc-900 border border-zinc-800 rounded-md px-3 py-2 text-sm text-white w-64 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 outline-none" 
                placeholder="Örn: HackerNews"
                required
              />
            </div>
            <div>
              <label className="block text-xs text-zinc-500 font-bold uppercase tracking-wider mb-1">Slug (Benzersiz Tanımlayıcı)</label>
              <input 
                type="text" 
                value={newTopic.slug} 
                onChange={e => setNewTopic({...newTopic, slug: e.target.value})}
                className="bg-zinc-900 border border-zinc-800 rounded-md px-3 py-2 text-sm text-white w-64 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 outline-none" 
                placeholder="Örn: hackernews"
                required
              />
            </div>
            <button type="submit" disabled={addMutation.isPending} className="bg-blue-600 text-white font-bold text-sm px-6 py-2.5 rounded-lg hover:bg-blue-500 transition-all">
              {addMutation.isPending ? 'Ekleniyor...' : 'Kaydet'}
            </button>
            <button type="button" onClick={() => setIsAdding(false)} className="border border-zinc-800 text-zinc-400 px-6 py-2.5 rounded-lg hover:border-zinc-700 hover:text-white transition-all">
              İptal
            </button>
          </form>
        </div>
      )}

      <div className="bg-zinc-950 border border-zinc-800 rounded-xl overflow-hidden shadow-lg shadow-black/40">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-zinc-900/60 border-b border-zinc-800 text-zinc-500">
                <th className="px-6 py-4 text-xs font-extrabold uppercase tracking-wider">Source Name</th>
                <th className="px-6 py-4 text-xs font-extrabold uppercase tracking-wider">Slug</th>
                <th className="px-6 py-4 text-xs font-extrabold uppercase tracking-wider text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-900">
              {isLoading ? (
                <tr><td colSpan="3" className="px-6 py-6 text-center text-zinc-500">Yükleniyor...</td></tr>
              ) : topics?.length > 0 ? (
                topics.map(topic => (
                  <tr key={topic.id} className="hover:bg-zinc-900/40 transition-colors group">
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-3">
                        <div className="w-2 h-2 rounded-full bg-emerald-500"></div>
                        <span className="text-sm font-bold text-white">{topic.name}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <span className="text-xs text-zinc-400 font-mono">{topic.slug}</span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex justify-end gap-2 opacity-70 group-hover:opacity-100 transition-opacity">
                        <button onClick={() => handleDelete(topic.id)} disabled={deleteMutation.isPending} className="p-1.5 bg-zinc-900 text-zinc-500 hover:bg-red-500/20 hover:text-red-400 rounded transition-colors" title="Delete">
                          <span className="material-symbols-outlined text-[18px]">delete</span>
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr><td colSpan="3" className="px-6 py-6 text-center text-zinc-500">Kaynak bulunamadı.</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </AdminLayout>
  );
}
