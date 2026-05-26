import React from 'react';
import { Link } from 'react-router-dom';
import { formatRelativeDate } from '../../utils/formatDate';
import './ArticleCard.css';

/**
 * Terminal Monolith article card — matches the Stitch "HaberBot - Keşfet" design.
 * Opaque #0A0A0A surface, #222222 border, JetBrains Mono labels, Geist headings.
 * Wide (featured) card at index 3 includes an image on the left side.
 */
const stripHtmlTags = (str) => {
  if (!str) return '';
  return str.replace(/<\/?[^>]+(>|$)/g, "").trim();
};

const ArticleCard = React.memo(function ArticleCard({ article, index = 0 }) {
  const {
    id,
    title_tr: titleTr,
    summary_tr: summaryTr,
    source,
    original_url: originalUrl,
    created_at: createdAt,
    image_url: imageUrl,
  } = article;

  const staggerClass = index < 6 ? `stagger-${index + 1}` : '';
  const rawSummary = summaryTr || article.summary || article.original_content || 'Özet hazırlanıyor...';
  const cleanSummary = stripHtmlTags(rawSummary);
  const displayTitle = titleTr || article.title;
  const displayDate = formatRelativeDate(createdAt);

  // Generate tags based on source
  const tags = source === 'hackernews'
    ? ['#Tech', '#HN']
    : source === 'reddit'
      ? ['#Reddit', '#Dev']
      : ['#Haber', '#Gündem'];

  // Index 3 is the wide/featured card (matches Stitch layout: card 4 spans 2 cols)
  const isWide = index === 3;
  const fallbackImage = "https://lh3.googleusercontent.com/aida-public/AB6AXuAkcc4ljifyLm60qOJU9nPdkSupP_lThlvDm-3PYXSuYnOAIxBrt_nukNvSRCfubrXcuhudGmrl1wl5XkItdKksmcKfoURVQoPp4OhrSNSOsTcRUTKlvB6Q6YIsmYBKHisQTS_Mvt2BRPC6HkjbUEYSEfchTrhdhKIQNM0ebjTlClXLobg7tNjt3XSF8pZpVsCwl3V6tZ256t_UoPSnv9gDJSAr2tJGGIGlEM3UvSpxDcmMtIlVK9FamSkNoefAeKIhjS7NnLSvuSk";
  const displayImage = imageUrl || fallbackImage;

  if (isWide) {
    return (
      <Link
        to={`/haber/${id}`}
        className={`article-card article-card--wide ${staggerClass}`}
        style={{ animationDelay: `${index * 0.04}s` }}
      >
        <div className="article-card__inner">
          {/* Image on the left — exactly like Stitch Card 4 */}
          <div className="article-card__image-wrap">
            <img
              src={displayImage}
              alt={displayTitle}
              className="article-card__image"
            />
          </div>

          {/* Content on the right */}
          <div className="article-card__body">
            <div className="article-card__meta">
              <span className="article-card__badge">{source || 'KAYNAK'}</span>
              <span className="article-card__date">{displayDate}</span>
            </div>
            <h2 className="article-card__title">{displayTitle}</h2>
            <p className="article-card__summary">{cleanSummary}</p>
            <div className="article-card__footer">
              {tags.map((tag, idx) => (
                <span key={idx} className="article-card__tag">{tag}</span>
              ))}
            </div>
          </div>
        </div>
      </Link>
    );
  }

  // Standard card (no image)
  return (
    <Link
      to={`/haber/${id}`}
      className={`article-card ${staggerClass}`}
      style={{ animationDelay: `${index * 0.04}s` }}
    >
      <div className="article-card__meta">
        <span className="article-card__badge">{source || 'KAYNAK'}</span>
        <span className="article-card__date">{displayDate}</span>
      </div>
      <h2 className="article-card__title">{displayTitle}</h2>
      <p className="article-card__summary">{cleanSummary}</p>
      <div className="article-card__footer">
        {tags.map((tag, idx) => (
          <span key={idx} className="article-card__tag">{tag}</span>
        ))}
      </div>
    </Link>
  );
});

export default ArticleCard;
