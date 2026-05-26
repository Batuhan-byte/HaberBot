import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useTopicArticles } from '../../hooks/useTopicArticles';
import { useTopics } from '../../hooks/useTopics';
import ArticleCard from '../../components/ArticleCard/ArticleCard';
import Skeleton from '../../components/Skeleton/Skeleton';
import Pagination from '../../components/Pagination/Pagination';
import ErrorState from '../../components/ErrorState/ErrorState';
import './TopicPage.css';

const TopicPage = () => {
  const { slug } = useParams();
  const [page, setPage] = useState(1);

  // Reset page when slug changes
  useEffect(() => {
    setPage(1);
  }, [slug]);

  const { data: topicsData } = useTopics();
  const { data: articlesData, isLoading, isError, refetch } = useTopicArticles(slug, page, 9);

  const topics = topicsData?.topics || [];
  const currentTopic = topics.find(t => t.slug === slug);
  const topicName = currentTopic ? currentTopic.name : 'Konu';

  const articles = articlesData?.articles || [];
  const totalArticles = articlesData?.total || 0;
  const limit = 9;
  const totalPages = Math.ceil(totalArticles / limit) || 1;

  const handlePageChange = (newPage) => {
    setPage(newPage);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  return (
    <div className="topic-page">
      <header className="topic-header">
        <h1 className="topic-title">
          <span className="gradient-text">{topicName}</span> Haberleri
        </h1>
        {currentTopic?.keywords && (
          <div className="topic-keywords">
            {currentTopic.keywords.map((kw, i) => (
              <span key={i} className="keyword-pill">#{kw}</span>
            ))}
          </div>
        )}
      </header>

      <section className="articles-section">
        {isError ? (
          <ErrorState message="Haberler yüklenirken bir sorun oluştu." onRetry={refetch} />
        ) : isLoading ? (
          <div className="articles-grid">
            {Array.from({ length: 6 }).map((_, i) => (
              <Skeleton key={i} />
            ))}
          </div>
        ) : articles.length === 0 ? (
          <div className="no-articles">
            <span className="no-articles-icon">📭</span>
            <h3>Haber Bulunamadı</h3>
            <p>Bu kategoriye ait haber henüz bulunmamaktadır.</p>
          </div>
        ) : (
          <>
            <div className="articles-grid">
              {articles.map((article, index) => (
                <ArticleCard 
                  key={article.id} 
                  article={article} 
                  index={index} 
                />
              ))}
            </div>
            {totalPages > 1 && (
              <div className="pagination-wrapper">
                <Pagination 
                  currentPage={page} 
                  totalPages={totalPages} 
                  onPageChange={handlePageChange} 
                />
              </div>
            )}
          </>
        )}
      </section>
    </div>
  );
};

export default TopicPage;
