import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { api } from '../../services/api';
import './SearchDropdown.css';

const FALLBACK_IMAGE = "https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?auto=format&fit=crop&w=600&q=80";

export default function SearchDropdown({ query, onSelect }) {
  const { data, isLoading } = useQuery({
    queryKey: ['search-suggestions', query],
    queryFn: () => api.searchArticles(query, '', 1, 4), // Fetch top 4 suggestions
    enabled: query.length >= 2,
    staleTime: 30 * 1000 // Cache for 30s
  });

  if (query.length < 2) return null;

  const articles = data?.articles || [];

  return (
    <div className="hb-search-dropdown-panel premium-glass-card">
      {isLoading ? (
        <div className="hb-dropdown-loading">
          <div className="hb-dropdown-spinner"></div>
          <span>Öneriler aranıyor...</span>
        </div>
      ) : articles.length === 0 ? (
        <div className="hb-dropdown-empty">
          <span className="material-symbols-outlined text-[20px] text-gray-500">search_off</span>
          <span>Eşleşen haber bulunamadı.</span>
        </div>
      ) : (
        <div className="hb-dropdown-list" role="listbox">
          <div className="hb-dropdown-header">
            <span className="hb-dropdown-glow-dot"></span>
            <span>ÖNERİLEN HABERLER</span>
          </div>
          {articles.map((article, index) => {
            const displayTitle = article.title_tr || article.title;
            const style = { animationDelay: `${index * 0.05}s` };
            return (
              <div 
                key={article.id} 
                className="hb-dropdown-item interactive-card"
                onClick={() => onSelect(article)}
                style={style}
                role="option"
                aria-selected="false"
              >
                <div className="hb-dropdown-thumb-wrap">
                  <img 
                    src={article.image_url || FALLBACK_IMAGE} 
                    alt={displayTitle} 
                    className="hb-dropdown-thumb" 
                    onError={(e) => { e.target.src = FALLBACK_IMAGE; }}
                  />
                </div>
                <div className="hb-dropdown-item-body">
                  <span className="hb-dropdown-item-title line-clamp-2">{displayTitle}</span>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}