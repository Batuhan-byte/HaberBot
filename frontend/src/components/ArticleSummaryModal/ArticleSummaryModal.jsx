import React, { useEffect, useRef } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeSanitize from 'rehype-sanitize';
import { useSummary } from '../../hooks/useSummary';
import SummarySkeleton from './SummarySkeleton';
import './ArticleSummaryModal.css';

/**
 * Premium AI Summary Modal with backdrop blur, custom skeleton loader,
 * client-side caching via TanStack Query and safe Markdown rendering.
 */
export default function ArticleSummaryModal({ articleId, open, onClose }) {
  const modalRef = useRef(null);

  // Fetch summary only when modal is open and article ID is present
  const { data, isLoading, error, refetch } = useSummary(articleId, open);

  useEffect(() => {
    if (!open) return;

    const handleKeyDown = (e) => {
      if (e.key === 'Escape') {
        onClose();
      }
    };

    // Lock background scrolling
    document.body.style.overflow = 'hidden';
    window.addEventListener('keydown', handleKeyDown);

    return () => {
      document.body.style.overflow = '';
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [open, onClose]);

  if (!open) return null;

  const handleOverlayClick = (e) => {
    if (modalRef.current && !modalRef.current.contains(e.target)) {
      onClose();
    }
  };

  return (
    <div 
      className="summary-modal-overlay" 
      onClick={handleOverlayClick}
      role="dialog"
      aria-modal="true"
      aria-labelledby="summary-modal-title"
    >
      <div className="summary-modal-container" ref={modalRef}>
        <header className="summary-modal-header">
          <h2 id="summary-modal-title" className="summary-modal-title">
            <span className="material-symbols-outlined summary-modal-title-sparkle" aria-hidden="true">auto_awesome</span>
            <span>HaberBot Yapay Zeka Özeti</span>
          </h2>
          <button 
            className="summary-modal-close-btn" 
            onClick={onClose}
            aria-label="Kapat"
            id="summary-modal-close-button"
          >
            <span className="material-symbols-outlined" aria-hidden="true">close</span>
          </button>
        </header>

        <div className="summary-modal-body custom-scrollbar">
          {isLoading ? (
            <SummarySkeleton />
          ) : error ? (
            <div className="summary-modal-error">
              <span className="material-symbols-outlined summary-modal-error-icon" aria-hidden="true">error</span>
              <p className="font-semibold">Özet Yüklenemedi</p>
              <p className="summary-modal-error-text">
                Yapay zeka özeti alınırken veya oluşturulurken bir hata oluştu. Lütfen tekrar deneyin.
              </p>
              <button className="summary-modal-error-retry" onClick={() => refetch()}>
                Tekrar Dene
              </button>
            </div>
          ) : data?.summary ? (
            <article className="summary-markdown-content animate-fade-in">
              <ReactMarkdown 
                remarkPlugins={[remarkGfm]} 
                rehypePlugins={[rehypeSanitize]}
              >
                {data.summary}
              </ReactMarkdown>
            </article>
          ) : (
            <div className="summary-modal-error">
              <span className="material-symbols-outlined summary-modal-error-icon" aria-hidden="true">warning</span>
              <p className="font-semibold">Özet Bulunamadı</p>
              <p className="summary-modal-error-text">
                Bu haber için herhangi bir yapay zeka özeti üretilemedi.
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
